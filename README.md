# VaspSuperCell
The usage of this tool can be obtained by running:
```console
  go run VaspSuperCell.go -h
```
which gives
```console
  -s string
    	POSCAR or CONTCAR (default "POSCAR")
  -v1 int
    	# of extension along vector 1 (default 1)
  -v2 int
    	# of extension along vector 2 (default 1)
  -v3 int
    	# of extension along vector 3 (default 1)
```
For example, if you want to generate a 2x3x4 supercell using the POSCAR in the current directory, please run 
```console
go run VaspSuperCell.go -v1 2 -v2 3 -v3 4 ./POSCAR
```
