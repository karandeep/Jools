#!/bin/bash
i="0"
while [ $i -lt 1 ]
do
  echo "Another run"
  perl -pne 'BEGIN {undef $/} s/([0-9]+,"[A-Z0-9]+","[a-zA-Z0-9 ]+",\d+,\d*\.?\d*,\d*\.?\d*,\d*\.?\d*,-?\d*\.?\d*),"?([0-9. x;]*)"?,"?([0-9.;]*)"?,([0-9]*,[0-9]*),"?([0-9. x;]*)"?,"?([0-9.;]*)"?,([0-9]*,[0-9]*),"?([0-9.;]*)"?,"?([0-9.;]*)"?\n,,,,,,,([0-9. x]*),(\d*\.?\d*),\d*,\d*,(\d*\.?\d*),(\d*\.?\d*),\d*,\d*,(\d*\.?\d*),(\d*\.?\d*).*\n/$1,"$2;$10;","$3;$11;",$4,"$5;$12;","$6;$13;",$7,"$8;$14;","$9;$15;"\n/g;' designs.csv > designs1.csv

  diff -u designs.csv designs1.csv > designs.diff

  if [ -s designs.diff ]
  then
    cp designs1.csv designs.csv
  else
    i=$[$i+1]
  fi
done
