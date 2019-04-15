package protein

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

const (
	MODEL_DAYHOFF = iota
	MODEL_JTT
	MODEL_MTREV
	MODEL_LG
	MODEL_WAG

	PROT_DIST_MAX = 20.00
	BL_MIN        = 1.e-08
	BL_MAX        = 100.0
	DBL_MIN       = 2.2250738585072014e-308
)

var DBL_EPSILON float64 = math.Nextafter(1, 2) - 1

type ProtModel struct {
	pi         []float64  // aa frequency
	mat        *mat.Dense // substitution matrix
	mr         float64    //MeanRate
	eigen      *mat.Eigen // Eigen Values/vectors
	leigenvect *mat.Dense // Left Eigen Vector (Inv of Eigen Vector)
	reigenvect *mat.Dense // Right Eigen Vector
	eval       []float64  // Eigen values
	ns         int        // Number of states in the model
	pij        *mat.Dense // Matrix of Pij
	alpha      float64    // Alpha
	stepsize   int
	usegamma   bool
}

// Initialize a new protein model, given the name of the model as const int:
// MODEL_DAYHOFF, MODEL_JTT, MODEL_MTREV, MODEL_LG or MODEL_WAG
func NewProtModel(model int, usegamma bool, alpha float64) (*ProtModel, error) {
	var m *mat.Dense
	var pi []float64
	switch model {
	case MODEL_DAYHOFF:
		m, pi = DayoffMats()
	case MODEL_JTT:
		m, pi = JTTMats()
	case MODEL_MTREV:
		m, pi = MtREVMats()
	case MODEL_LG:
		m, pi = LGMats()
	case MODEL_WAG:
		m, pi = WAGMats()
	default:
		return nil, fmt.Errorf("This protein model is not implemented")
	}
	return &ProtModel{
		pi,
		m,
		-1.0,
		nil,
		nil,
		nil,
		nil,
		len(pi),
		mat.NewDense(len(pi), len(pi), nil),
		alpha,
		1,
		usegamma,
	}, nil
}

// Returns code of the model
// If the model does not exist, returns -1
func ModelStringToInt(model string) int {
	switch model {
	case "dayoff":
		return MODEL_DAYHOFF
	case "jtt":
		return MODEL_JTT
	case "mtrev":
		return MODEL_MTREV
	case "lg":
		return MODEL_LG
	case "wag":
		return MODEL_WAG
	default:
		return -1
	}
}

// Initialize model with given aa frequencies
// If aafreqs is null, then uses model frequencies
func (model *ProtModel) InitModel(aafreqs []float64) error {
	var i, j int
	var sum float64
	var ok bool

	ns := 20

	if model.mat == nil || model.pi == nil {
		return fmt.Errorf("Matrices have not been initialized")
	}

	if aafreqs != nil && len(aafreqs) != 20 {
		return fmt.Errorf("aa frequency array does not have a length of 20")
	}
	model.pi = aafreqs

	/* multiply the nth col of Q by the nth term of pi/100 just as in PAML */
	model.mat.Apply(func(i, j int, v float64) float64 { return v * model.pi[j] / 100.0 }, model.mat)

	/* compute diagonal terms of Q and mean rate mr = l/t */
	model.mr = .0
	for i = 0; i < ns; i++ {
		sum = .0
		for j = 0; j < ns; j++ {
			sum += model.mat.At(i, j)
		}
		model.mat.Set(i, i, -sum)
		model.mr += model.pi[i] * sum
	}

	/* scale instantaneous rate matrix so that mu=1 */
	model.mat.Apply(func(i, j int, v float64) float64 { return v / model.mr }, model.mat)

	model.eigen = &mat.Eigen{}
	if ok = model.eigen.Factorize(model.mat, mat.EigenRight); !ok {
		return fmt.Errorf("Problem during matrix decomposition")
	}
	model.reigenvect = mat.NewDense(ns, ns, nil)
	model.leigenvect = mat.NewDense(20, 20, nil)
	u := model.eigen.VectorsTo(nil)
	model.eval = make([]float64, ns)
	for i, b := range model.eigen.Values(nil) {
		model.eval[i] = real(b)
	}
	model.reigenvect.Apply(func(i, j int, val float64) float64 { return real(u.At(i, j)) }, model.reigenvect)
	model.leigenvect.Inverse(model.reigenvect)

	return nil
}

func (model *ProtModel) Ns() int {
	return model.ns
}

func (model *ProtModel) PMat(l float64) {
	if l < BL_MIN {
		model.pMatZeroBrLen()
	} else {
		model.pMatEmpirical(l)
	}
}

func (model *ProtModel) pMatZeroBrLen() {
	model.pij.Apply(func(i, j int, v float64) float64 {
		if i == j {
			return 1.0
		}
		return 0.0
	}, model.pij)
}

/********************************************************************/

/* Computes the substitution probability matrix
 * from the initial substitution rate matrix and frequency vector
 * and one specific branch length
 *
 * input : l , branch length
 * input : mod , choosen model parameters, qmat and pi
 * ouput : Pij , substitution probability matrix
 *
 * matrix P(l) is computed as follows :
 * P(l) = exp(Q*t) , where :
 *
 *   Q = substitution rate matrix = Vr*D*inverse(Vr) , where :
 *
 *     Vr = right eigenvector matrix for Q
 *     D  = diagonal matrix of eigenvalues for Q
 *
 *   t = time interval = l / mr , where :
 *
 *     mr = mean rate = branch length/time interval
 *        = sum(i)(pi[i]*p(i->j)) , where :
 *
 *       pi = state frequency vector
 *       p(i->j) = subst. probability from i to a different state
 *               = -Q[ii] , as sum(j)(Q[ij]) +Q[ii] = 0
 *
 * the Taylor development of exp(Q*t) gives :
 * P(l) = Vr*exp(D*t)        *inverse(Vr)
 *      = Vr*pow(exp(D/mr),l)*inverse(Vr)
 *
 * for performance we compute only once the following matrices :
 * Vr, inverse(Vr), exp(D/mr)
 * thus each time we compute P(l) we only have to :
 * make 20 times the operation pow()
 * make 2 20x20 matrix multiplications, that is :
 *   16000 = 2x20x20x20 times the operation *
 *   16000 = 2x20x20x20 times the operation +
 *   which can be reduced to (the central matrix being diagonal) :
 *   8400 = 20x20 + 20x20x20 times the operation *
 *   8000 = 20x20x20 times the operation + */
func (model *ProtModel) pMatEmpirical(len float64) {
	var i, k int
	var U, V *mat.Dense
	var R []float64
	var expt []float64
	var uexpt *mat.Dense
	var tmp float64

	U = model.reigenvect //mod->eigen->r_e_vect;
	R = model.eval       //mod->eigen->e_val;// To take only real part from that vector /* eigen value matrix */
	V = model.leigenvect
	expt = make([]float64, model.ns)              //model.eigen.Values(nil) // To take only imaginary part from that vector
	uexpt = mat.NewDense(model.ns, model.ns, nil) //model.eigen.Vectors() //  don't know yet how to handle that // mod->eigen->r_e_vect_im;

	model.pij.Apply(func(i, j int, v float64) float64 { return .0 }, model.pij)
	tmp = .0

	for k = 0; k < model.ns; k++ {
		expt[k] = R[k]
	}

	if model.usegamma && (math.Abs(model.alpha) > DBL_EPSILON) {
		// compute pow (alpha / (alpha - e_val[i] * l), alpha)
		for i = 0; i < model.ns; i++ {
			tmp = model.alpha / (model.alpha - (R[i] * len))
			expt[i] = math.Pow(tmp, model.alpha)
		}
	} else {
		for i = 0; i < model.ns; i++ {
			expt[i] = float64(math.Exp(R[i] * len))
		}
	}

	// multiply Vr* pow (alpha / (alpha - e_val[i] * l), alpha) *Vi into Pij
	uexpt.Apply(func(i, j int, v float64) float64 {
		return U.At(i, j) * expt[j]
	}, uexpt)
	model.pij.Apply(func(i, j int, v float64) float64 {
		for k = 0; k < model.ns; k++ {
			v += uexpt.At(i, k) * V.At(k, j)
		}
		if v < DBL_MIN {
			v = DBL_MIN
		}
		return v

	}, model.pij)
}
