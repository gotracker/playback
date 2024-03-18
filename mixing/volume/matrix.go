package volume

import "math"

// Matrix is an array of Volumes
type Matrix struct {
	StaticMatrix
	Channels int
}

const (
	StereoCoeff = math.Sqrt2 / 2.0
)

// Apply takes a volume matrix and multiplies it by incoming volumes
func (m Matrix) ApplyToMatrix(mtx Matrix) Matrix {
	if mtx.Channels == 0 {
		return m
	}

	if m.Channels == mtx.Channels {
		// simple straight-through
		for i := 0; i < m.Channels; i++ {
			m.StaticMatrix[i] = mtx.StaticMatrix[i].ApplySingle(m.StaticMatrix[i])
		}
		return m
	}

	// more complex applications follow...

	if mtx.Channels == 1 {
		// right (mtx) is mono, so just do direct mono application
		return m.Apply(mtx.StaticMatrix[0])
	}

	// NOTE: recursive
	return m.ApplyToMatrix(mtx.ToChannels(m.Channels))
}

func (m Matrix) Apply(vol Volume) Matrix {
	for i := 0; i < m.Channels; i++ {
		m.StaticMatrix[i] = vol.ApplySingle(m.StaticMatrix[i])
	}
	return m
}

func (m *Matrix) Accumulate(in Matrix) {
	if m.Channels == 0 {
		*m = in
		return
	}

	dry := in.ToChannels(m.Channels)
	for i := 0; i < m.Channels; i++ {
		m.StaticMatrix[i] += dry.StaticMatrix[i]
	}
}

func (m *Matrix) Assign(channels int, data []Volume) {
	m.Channels = channels
	for i := 0; i < channels; i++ {
		m.StaticMatrix[i] = data[i]
	}
}

func (m Matrix) ToChannels(channels int) Matrix {
	if m.Channels == channels {
		return m
	}

	switch channels {
	case 1:
		return m.AsMono()
	case 2:
		return m.AsStereo()
	case 4:
		return m.AsQuad()
	default:
		return Matrix{}
	}
}

// Sum sums all the elements of the Matrix and returns the resulting Volume
func (m Matrix) Sum() Volume {
	switch m.Channels {
	case 0:
		return 0
	case 1:
		return m.StaticMatrix[0]
	case 2:
		return (m.StaticMatrix[0] + m.StaticMatrix[1]) / StereoCoeff
	case 4:
		return (m.StaticMatrix[0] + m.StaticMatrix[1] + m.StaticMatrix[2] + m.StaticMatrix[3]) / StereoCoeff
	default:
		c := Volume(1 / float64(m.Channels))

		var v Volume
		for i := 0; i < m.Channels; i++ {
			v += m.StaticMatrix[i] * c
		}
		return v
	}
}

func (m *Matrix) Set(ch int, vol Volume) {
	m.StaticMatrix[ch] = vol
}

func (m Matrix) Get(ch int) Volume {
	return m.StaticMatrix[ch]
}

func (m Matrix) AsMono() Matrix {
	switch m.Channels {
	case 0:
		return Matrix{}
	case 1:
		return m
	default:
		return Matrix{
			StaticMatrix: StaticMatrix{m.Sum()},
			Channels:     1,
		}
	}
}

func (m Matrix) AsStereo() Matrix {
	switch m.Channels {
	case 0:
		return Matrix{}
	case 1:
		return Matrix{
			StaticMatrix: StaticMatrix{m.StaticMatrix[0] * StereoCoeff, m.StaticMatrix[0] * StereoCoeff},
			Channels:     2,
		}
	case 2:
		return m
	case 4:
		return Matrix{
			StaticMatrix: StaticMatrix{(m.StaticMatrix[0] + m.StaticMatrix[2]) / 2.0, (m.StaticMatrix[1] + m.StaticMatrix[3]) / 2.0},
			Channels:     2,
		}
	default:
		return Matrix{}
	}
}

func (m Matrix) AsQuad() Matrix {
	switch m.Channels {
	case 0:
		return Matrix{}
	case 1:
		return Matrix{
			StaticMatrix: StaticMatrix{m.StaticMatrix[0] * StereoCoeff, m.StaticMatrix[0] * StereoCoeff, m.StaticMatrix[0] * StereoCoeff, m.StaticMatrix[0] * StereoCoeff},
			Channels:     4,
		}
	case 2:
		return Matrix{
			StaticMatrix: StaticMatrix{m.StaticMatrix[0], m.StaticMatrix[1], m.StaticMatrix[0], m.StaticMatrix[1]},
			Channels:     4,
		}
	case 4:
		return m
	default:
		return Matrix{}
	}
}

func (m Matrix) Lerp(other Matrix, t float32) Matrix {
	if other.Channels == 0 || t <= 0 {
		return m
	}

	out := other.ToChannels(m.Channels)

	// lerp between m and v
	for c := 0; c < m.Channels; c++ {
		a := m.StaticMatrix[c]
		b := out.StaticMatrix[c]
		out.StaticMatrix[c] = a + Volume(t)*(b-a)
	}
	return out
}
