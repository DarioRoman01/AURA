path=`go list -f '{{.Target}}'`
set PATH=%PATH%;C:path
go install && (
    echo aura was installed successfuly
) || (
    echo installation fail
)