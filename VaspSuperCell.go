package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadLines(textfile string) (lines []string, err error) {
	f, err := os.Open(textfile)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return
}

func StrsToInts(strs []string) []int {
	var is = []int{}
	for _, x := range strs {
		t, _ := strconv.Atoi(x)
		is = append(is, t)
	}

	return is
}

func StrsToFloat64s(strs []string) []float64 {
	var fs = []float64{}
	for _, x := range strs {
		t, _ := strconv.ParseFloat(x, 64)
		fs = append(fs, t)
	}
	return fs
}

func StrToFloat64(str string) float64 {
	t, _ := strconv.ParseFloat(strings.TrimSpace(str), 64)
	return t
}

func LineToFloat64s(str string) []float64 {
	sl := strings.Fields(str)
	var fs = []float64{}

	for _, x := range sl {
		t, _ := strconv.ParseFloat(x, 64)
		fs = append(fs, t)
	}
	return fs
}

func SumSliceInt(s []int) int {
	var sum int = 0
	for _, x := range s {
		sum += x
	}
	return sum
}

type Position struct {
	x float64
	y float64
	z float64
}

type Atom struct {
	symbol  string
	x, y, z float64
}

type POSCAR struct {
	comment    string
	latt_const float64
	latt_a     []float64
	latt_b     []float64
	latt_c     []float64
	atoms      []string
	natoms     []int
	ntotatom   int
	isfrac     bool
	pos        []Position
}

func extend(p POSCAR, n1 int, n2 int, n3 int) (sp POSCAR) {

	sp.comment = "This is a supercell created with VaspSuperCell.go"
	sp.latt_const = p.latt_const

	for i := 0; i < 3; i++ {
		sp.latt_a = append(sp.latt_a, p.latt_a[i]*float64(n1))
		sp.latt_b = append(sp.latt_b, p.latt_b[i]*float64(n2))
		sp.latt_c = append(sp.latt_c, p.latt_c[i]*float64(n3))
	}

	n123 := n1 * n2 * n3

	istart := 0
	iend := 0
	for i := 0; i < len(p.atoms); i++ {
		iend += p.natoms[i]

		sp.atoms = append(sp.atoms, p.atoms[i])
		sp.natoms = append(sp.natoms, p.natoms[i]*n123)

		atoms := p.pos[istart:iend]

		for j1 := 0; j1 < n1; j1++ {
			for j2 := 0; j2 < n2; j2++ {
				for j3 := 0; j3 < n3; j3++ {

					for _, at := range atoms {
						var t Position
						t.x = (at.x + float64(j1)) / float64(n1)
						t.y = (at.y + float64(j2)) / float64(n2)
						t.z = (at.z + float64(j3)) / float64(n3)
						sp.pos = append(sp.pos, t)
					}

				}
			}
		}

		istart = iend
	}

	sp.ntotatom = SumSliceInt(sp.natoms)
	sp.isfrac = p.isfrac
	
	return
}

func display(p POSCAR) {
	fmt.Println(p.comment)

	fmt.Printf("%s\n", "Lattice Constant:")
	fmt.Printf("%20.12f\n", p.latt_const)

	fmt.Printf("%s\n", "Vector a:")
	fmt.Printf("%20.12f\n", p.latt_a)

	fmt.Printf("%s\n", "Vector b:")
	fmt.Printf("%20.12f\n", p.latt_b)

	fmt.Printf("%s\n", "Vector c:")
	fmt.Printf("%20.12f\n", p.latt_c)

	fmt.Println(p.atoms)

	fmt.Printf("%6d\n", p.natoms)

	fmt.Println(p.isfrac)

	for i, _ := range p.pos {
		fmt.Printf("%20.12f\n", p.pos[i])
	}
}

func parse(poscarfile string) (poscar POSCAR) {
	lines, _ := ReadLines(poscarfile)

	poscar.comment = lines[0]
	poscar.latt_const = StrToFloat64(lines[1])
	poscar.latt_a = LineToFloat64s(lines[2])
	poscar.latt_b = LineToFloat64s(lines[3])
	poscar.latt_c = LineToFloat64s(lines[4])

	poscar.atoms = strings.Fields(lines[5])
	poscar.natoms = StrsToInts(strings.Fields(lines[6]))

	poscar.ntotatom = SumSliceInt(poscar.natoms)

	if strings.TrimSpace(lines[7])[0] == 'D' {
		poscar.isfrac = true
	} else {
		poscar.isfrac = false
	}

	base := 8
	for i := 0; i < poscar.ntotatom; i++ {
		tpos := LineToFloat64s(lines[base+i])
		var pos Position
		pos.x = tpos[0]
		pos.y = tpos[1]
		pos.z = tpos[2]

		poscar.pos = append(poscar.pos, pos)

	}
	return
}

func get_opts() (p_n1, p_n2, p_n3 *int, p_poscarfile *string) {

	p_n1 = flag.Int("v1", 1, "# of extension along vector 1")
	p_n2 = flag.Int("v2", 1, "# of extension along vector 2")
	p_n3 = flag.Int("v3", 1, "# of extension along vector 3")
	p_poscarfile = flag.String("s", "POSCAR", "POSCAR or CONTCAR")

	flag.Parse()

	return
}

func output_vasp(p POSCAR) {
     fmt.Printf("%s\n", p.comment)
     fmt.Printf("%20.16f\n", p.latt_const)

     var v = p.latt_a
     fmt.Printf("%22.16f%22.16f%22.16f\n", v[0], v[1], v[2])

     v = p.latt_b
     fmt.Printf("%22.16f%22.16f%22.16f\n", v[0], v[1], v[2])

     v = p.latt_c
     fmt.Printf("%22.16f%22.16f%22.16f\n", v[0], v[1], v[2])

     for _,at := range p.atoms {
     	fmt.Printf("%5s", at)      
     }
     fmt.Printf("\n")
     
     for _,n := range p.natoms {
     	 fmt.Printf("%5d", n)
     }
     fmt.Printf("\n")

     if p.isfrac {
     	fmt.Printf("%s\n", "Direct")
     }

     for _, t := range p.pos {
     	 fmt.Printf("%22.16f%22.16f%22.16f\n", t.x,t.y,t.z)
     }

}

func main() {

	p_n1, p_n2, p_n3, p_poscarfile := get_opts()

	fmt.Println(*p_n1, *p_n2, *p_n3)

	poscar := parse(*p_poscarfile)

	display(poscar)

	sp := extend(poscar, *p_n1, *p_n2, *p_n3)
//	display(sp)
	output_vasp(sp)
}