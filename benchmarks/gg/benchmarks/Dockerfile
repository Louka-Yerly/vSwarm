# MIT License
#
# Copyright (c) 2021 Michal Baczun and EASE lab
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

FROM vhiveease/golang-builder AS builder

WORKDIR /app
# We use ARG and ENV to have a single Dockerfile for the all three programs.
# https://vsupalov.com/docker-build-pass-environment-variables/
ARG target_arg
ENV target=$target_arg
COPY . .
RUN go mod tidy
RUN go mod download && \
   CGO_ENABLED=0 GOOS=linux go build -v -o ./bin .

FROM vhiveease/runner
WORKDIR /app
ARG target_arg
ENV target=$target_arg
COPY ./setup_gg.sh .
RUN ./setup_gg.sh
COPY ./${target} ./${target}
ENV GG_LAMBDA_ROLE=arn:aws:iam::356764711652:role/gg
ENV AWS_REGION=us-west-1
ENV GG_MODELPATH=/app/gg/src/models/wrappers
COPY --from=builder /app/bin /app/exe
ENTRYPOINT [ "/app/exe" ]
