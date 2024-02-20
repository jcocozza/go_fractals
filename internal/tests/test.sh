# The following commands should all run successfully and produce content.
go run . ifs image -p examples/barnsley_fern_ifs.ifs --algo-d -n 13 -F test1
go run . ifs image -p examples/barnsley_fern_ifs.ifs --algo-p -n 67108864 -F test2
go run . ifs image -p examples/barnsley_fern_ifs.ifs --algo-p -n 67108864 --probabilities .1,.67,.115,.115 -F test3
go run . ifs image -p examples/maple/maple.ifs --algo-d -n 11 -F test4
go run . ifs image -p examples/maple/maple.ifs --algo-p -n 4194304 -F test5
go run . ifs evolve -p examples/maple/maple.ifs --algo-d -n 11 --fps 3 -F test6
go run . ifs evolve -p examples/maple/maple.ifs --threeDim -k 1 -T 50 -n 1000000 --algo-p -F test7
go run . julia -e "1/(z*z + .72i)" -F test8
go run . julia-evolve -e "1/(z*z + c)" -f 10 -P "0-0.63i" -n 100 -i "0-0.001i" -F test9
go run . julia-evolve -e "z*z + c" -P "-.5+0i" -n 10 -i ".0625+.0625i" -W 200 -H 200 --threeDim -F test10
go run . mandelbrot -e "z*z + c" -F test11
go run . mandelbrot -e "complex(math.Abs(real(z)),math.Abs(imag(z)))*complex(math.Abs(real(z)),math.Abs(imag(z))) + c" -F test12
go run . mandelbrot -e "complex(math.Abs(real(z)),math.Abs(imag(z)))*complex(math.Abs(real(z)),math.Abs(imag(z))) + c" -p "-1.75-0.025i" --color -z .08 -F test13