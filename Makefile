fmt:
	ls -d */ | grep -v vendor | xargs -I {} gofmt -s -w {}