cd pkg || exit; go mod tidy; cd ..;
cd services || exit;
for dir in gateway match profile rengine session
do
  cd $dir || exit; pwd; go mod tidy; cd ..;
done
