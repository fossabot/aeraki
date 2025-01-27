# Copyright Aeraki Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

BASEDIR=$(dirname "$0")
source $BASEDIR/../common_func.sh

kubectl create ns thrift
LabelIstioInjectLabel thrift
kubectl apply -f $BASEDIR/thrift-sample.yaml -n thrift
kubectl apply -f $BASEDIR/destinationrule.yaml -n thrift
kubectl apply -f $BASEDIR/virtualservice-traffic-splitting.yaml -n thrift
