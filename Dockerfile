FROM golang:1.3.3-onbuild
ENTRYPOINT [ "go-wrapper", "run" ]
CMD [ "" ]
