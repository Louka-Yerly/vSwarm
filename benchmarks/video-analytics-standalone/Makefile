# MIT License

# Copyright (c) 2024 EASE lab

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

ROOT = ../../

all-image: init-database video-analytics-standalone-python

video-analytics-standalone-python: docker/Dockerfile python/server.py python/imagenet_labels.txt python/requirements.txt videos/default.mp4
	DOCKER_BUILDKIT=1 docker build \
	--tag vhiveease/video-analytics-standalone-python:latest \
	--target videoAnalyticsStandalonePython \
	-f docker/Dockerfile \
	$(ROOT) --load

init-database: docker/Dockerfile init/init-database.go videos/default.mp4
	DOCKER_BUILDKIT=1 docker build \
	--tag vhiveease/video-analytics-standalone-init-database:latest \
	--target databaseInit \
	-f docker/Dockerfile \
	$(ROOT) --load

## Push images
push:
	docker push docker.io/vhiveease/video-analytics-standalone-python:latest 
	docker push docker.io/vhiveease/video-analytics-standalone-init-database:latest

## Pull images from docker hub
pull:
	docker pull docker.io/vhiveease/video-analytics-standalone-python:latest
	docker pull docker.io/vhiveease/video-analytics-standalone-init-database:latest
