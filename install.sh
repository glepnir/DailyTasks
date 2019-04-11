echo "changing directory to DailyTasks"
cd $GOPATH/src/github.com/taigacute/DailyTasks/config
echo "creating table"
cat schema.sql | sqlite3 tasks.db
cd $GOPATH/src/github.com/taigacute/DailyTasks
echo "building the go binary"
go build -o DailyTasks

echo "starting the binary"
./DailyTasks
