HOST=127.0.0.1
PORT=5555

echo "@ insert" 
curl -XPOST http://${HOST}:${PORT}/storage/test1 -d "Hello, hello world 1"
curl -XGET http://${HOST}:${PORT}/storage/test1 
echo

echo "@ update"
curl -XPOST http://${HOST}:${PORT}/storage/test1 -d "Hello, hello world 2"
curl -XGET http://${HOST}:${PORT}/storage/test1 
echo

echo "@ delete"
curl -XDELETE http://${HOST}:${PORT}/storage/test1
curl -XGET http://${HOST}:${PORT}/storage/test1 
echo