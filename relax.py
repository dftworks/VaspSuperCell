#!/usr/bin/env python
import sys
import math
import subprocess

MAXITER=20
TOLERANCE=1                             # in kB

def log(msg):
    f = open("run.log","wa")
    f.write(msg)
    f.close()

def readlines(filename):
    fil = open(filename)
    lines = fil.readlines()
    fil.close()
    return lines

def getpressure(lines):
    pline = ""
    for line in lines:        
        n = line.find("in kB")
        if n > 0:
            pline = line
    if pline == "":
        print "The information of pressure is not found"
    sl = pline.split()
    xx = float(sl[2])
    yy = float(sl[3])
    zz = float(sl[4])
    return [xx,yy,zz]

def getlattice(lines):
    latt = []
    sl = lines[2].split()
    latt.append([float(sl[0]), float(sl[1]), float(sl[2])])
    sl = lines[3].split()
    latt.append([float(sl[0]), float(sl[1]), float(sl[2])])
    sl = lines[4].split()
    latt.append([float(sl[0]), float(sl[1]), float(sl[2])])
    return latt

def runvasp(msg):
    print "run vasp %s" % (msg)
    subprocess.call('sh runvasp.sh "%s"' % msg, shell='True')

def buildposcar(relax_dir,newlatt):
    lines_con = readlines('CONTCAR')
    lines = []
        
    if relax_dir == "x" :
        for n in range(2):
            lines.append(lines_con[n])

        lines.append("%20.12f %20.12f %20.12f\n" % (newlatt,0.0,0.0))
        for line in lines_con[3:]:
            lines.append(line)
            
    elif relax_dir == "y" :
        for n in range(3):
            lines.append(lines_con[n])
        lines.append("%20.12f %20.12f %20.12f\n" % (0.0,newlatt,0.0))
        for line in lines_con[4:]:
            lines.append(line)
        
    elif relax_dir == "z" :
        for n in range(4):
            lines.append(lines_con[n])
        lines.append("%20.12f %20.12f %20.12f\n" % (0.0,0.0,newlatt))
        for line in lines_con[5:]:
            lines.append(line)
    elif relax_dir == "xy" :
        # the lattice vectors along x and y will be scaled by newlatt
        latt = getlattice(lines_con)
        zz = latt[2][2]
        latt[2][2] = zz * (float(lines_con[1])) / newlatt

        lines.append(lines_con[0])
        lines.append("%20.12f\n" % newlatt)

        for i in range(3):
            lines.append("%20.12f %20.12f %20.12f\n" % (latt[i][0], latt[i][1], latt[i][2]))
            
        for line in lines_con[5:]:
            lines.append(line)
    else:
        print "unimplemented lattice optimization in directions %s" % (relax_dir)

    f = open('POSCAR', 'w')
    for line in lines:
        f.write(line)
    f.close()

def pressure(relax_dir):
    lines_out = readlines('OUTCAR')
    prs = getpressure(lines_out)
    pp = 0.0
    if relax_dir == "x":
        pp = prs[0]
    elif relax_dir == "y":
        pp = prs[1]
    elif relax_dir == "z":
        pp = prs[2]
    elif relax_dir == "xy":
        pp = prs[0]
    else:
        print "unimplemented lattice optimization in directions %s" % (relax_dir)
    return pp

def lattice(relax_dir):
    lines_con = readlines('CONTCAR')
    latt = getlattice(lines_con)
    tt = 0.0
    if relax_dir == "x":
        tt = latt[0][0]
    elif relax_dir == "y":
        tt = latt[1][1]
    elif relax_dir == "z":
        tt = latt[2][2]
    elif relax_dir == "xy":
        tt = float(lines_con[1])
    else:
	print "unimplemented lattice optimization in directions %s" % (relax_dir)
    return tt

def calc_b(relax_dir,a):
    for i in range(1,100):
        if relax_dir in ("x","y","z"):
            xx = a-i*0.1
        elif relax_dir in ("xy","yz","xz"):
            xx = a-i*0.1
        else:
            pass
        
        buildposcar(relax_dir,xx)
        runvasp("b.%d" % (i))
        pres = pressure(relax_dir)
        if ( pres > 0 ):
            b = lattice(relax_dir)
            break

    return b

def repos_a(relax_dir,b):
    if relax_dir in ("x","y","z"):
        a = b + 0.1
    elif relax_dir in ("xy","yz","xz"):
        a = b + 0.1
    else:
        pass
    return a

def repos_b(relax_dir,a):
    if relax_dir in ("x","y","z"):
        b = a - 0.1
    elif relax_dir in ("xy","yz","xz"):
        b = a - 0.1
    else:
        pass
    return b

def calc_a(relax_dir,b):
    for i in range(1,100):
        if relax_dir in ("x","y","z"):
            xx = b+i*0.1
        elif relax_dir in ("xy","yz","xz"):
            xx = b+i*0.1
        else:
            pass
        
        buildposcar(relax_dir,xx)
        runvasp("a.%d" % (i))
        pres = pressure(relax_dir)
        if ( pres < 0 ):
            a = lattice(relax_dir)
            break
    return a

        
def findlimits(relax_dir):    
    runvasp("scf")
    pres = pressure(relax_dir)
    if pres < 0:
        a = lattice(relax_dir)
        b = calc_b(relax_dir,a)
        a = repos_a(relax_dir,b)
    else:
        b = lattice(relax_dir)
        a = calc_a(relax_dir,b)
        b = repos_b(relax_dir,a)
    print "found a,b: %20.12f, %20.12f" % (a,b)
    return [a,b]

def bisection(relax_dir,a,b):
    for m in range(1,MAXITER):
        c = (a+b)/2.0
        print "a,b,c = %20.12f,%20.12f,%20.12f" % (a,b,c)
        buildposcar(relax_dir,c)
        runvasp("bis.%d" %(m))
        pres=pressure(relax_dir)
        if abs(pres) < TOLERANCE:
            print "Relaxation completed in %d steps" % (m)
            break
        if pres < 0:
            a=c
            b=b
        else:
            b=c
            a=a
        
            
if __name__ == '__main__':
    relax_dir = sys.argv[1]
#    runvasp("scf")
    [a,b] = findlimits(relax_dir)
    bisection(relax_dir,a,b)
    
