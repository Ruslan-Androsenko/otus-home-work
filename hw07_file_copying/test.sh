#!/usr/bin/env bash
set -xeuo pipefail

go build -o go-cp

./go-cp -from testdata/input.txt -to out.txt
cmp out.txt testdata/out_offset0_limit0.txt

./go-cp -from testdata/input.txt -to out.txt -limit 10
cmp out.txt testdata/out_offset0_limit10.txt

./go-cp -from testdata/input.txt -to out.txt -limit 1000
cmp out.txt testdata/out_offset0_limit1000.txt

./go-cp -from testdata/input.txt -to out.txt -limit 10000
cmp out.txt testdata/out_offset0_limit10000.txt

./go-cp -from testdata/input.txt -to out.txt -offset 100 -limit 1000
cmp out.txt testdata/out_offset100_limit1000.txt

./go-cp -from testdata/input.txt -to out.txt -offset 6000 -limit 1000
cmp out.txt testdata/out_offset6000_limit1000.txt

# My tests begin
./go-cp -from testdata/input.txt -to out.txt -offset 100 -limit 2000
cmp out.txt testdata/out_offset100_limit2000.txt

./go-cp -from testdata/input.txt -to out.txt -offset 1000 -limit 2000
cmp out.txt testdata/out_offset1000_limit2000.txt

./go-cp -from testdata/input.txt -to out.txt -offset 2500 -limit 3050
cmp out.txt testdata/out_offset2500_limit3050.txt

./go-cp -from testdata/input.txt -to out.txt -offset 5000 -limit 2000
cmp out.txt testdata/out_offset5000_limit2000.txt

./go-cp -from testdata/input.txt -to out.txt -offset -6000 -limit 1000
cmp out.txt testdata/out_offset-6000_limit1000.txt

./go-cp -from testdata/input.txt -to out.txt -offset -200 -limit 1000
cmp out.txt testdata/out_offset-200_limit1000.txt
# My tests end

rm -f go-cp out.txt
echo "PASS"
