echo "changing directory to DailyTasks"
cd $GOPATH/src/github.com/taigacute/Tasks/config
echo "creating table"
cat schema.sql | sqlite3 tasks.db
echo "building the go binary"
go build -o Tasks

echo "starting the binary"
./Tasks
