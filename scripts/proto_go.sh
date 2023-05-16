#! /bin/sh

protoPath=$(pwd)/api/proto

protoFileList=""
for pkgDir in `ls $protoPath`;
do
    if [[ $pkgDir == "google" ]]; then
        continue
    fi
    for vDir in `ls $protoPath"/"$pkgDir`;
    do
        for file in `ls $protoPath"/"$pkgDir"/"$vDir`;
        do
            if [ -f $protoPath"/"$pkgDir"/"$vDir"/"$file ]; then
                protoFileList=$protoFileList" "$pkgDir"/"$vDir"/"$file
            fi
        done
    done
done
# echo "protoc -I api/proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $protoFileList"
protoExec=`protoc -I api/proto --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $protoFileList`