

old=\${$1}
new=$(eval echo \${$1})
sed -i "s#${old}#${new}#" $2


