#TAG=v0.0.2
# IMAGE=lawlerseth/hydro-scalar

# docker build -t $IMAGE:$TAG .

# docker run -it --entrypoint /bin/sh $IMAGE:$TAG 

# docker push $IMAGE:$TAG

# # test
# docker run --mount type=bind,src=/home/slawler/workbench/repos/go-wat/config,dst=/workspaces/config \
#      --mount type=bind,src=/home/slawler/workbench/repos/go-wat/test-data,dst=/workspaces/test-data \
#     $IMAGE:$TAG /bin/sh -c  "./main -config=config/event-settings.json"