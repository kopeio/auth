#!/usr/bin/env python

# Copyright 2017 The Kubernetes Authors.
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

# This python script helps sync godeps from the k8s repos into our git submodules
# It generates bash commands where changes are needed
# We can probably also use it for deps when the time comes!

import json
import sys
import subprocess
from pprint import pprint
from os.path import expanduser, join

kops_dir = expanduser('~/k8s/src/k8s.io/kops')
k8s_dir = expanduser('~/k8s/src/k8s.io/client-go')

with open(join(k8s_dir, 'Godeps/Godeps.json')) as data_file:
    godeps = json.load(data_file)

# pprint(godeps)

godep_map = {}

for godep in godeps['Deps']:
    # print("%s %s" % (godep['ImportPath'], godep['Rev']))
    godep_map[godep['ImportPath']] = godep['Rev']

# process = subprocess.Popen(['git', 'submodule', 'status'], stdout=subprocess.PIPE, cwd=kops_dir)
# submodule_status, err = process.communicate()
# for submodule_line in submodule_status.splitlines():
#   tokens = submodule_line.split()
#   dep = tokens[1]
#   dep = dep.replace('_vendor/', '')
#   sha = tokens[0]
#   sha = sha.replace('+', '')
#   godep_sha = godep_map.get(dep)
#   if not godep_sha:
#     for k in godep_map:
#       if k.startswith(dep):
#         godep_sha = godep_map[k]
#         break
#   if godep_sha:
#     if godep_sha != sha:
#       print("# update needed: %s vs %s" % (godep_sha, sha))
#       print("pushd _vendor/{dep}; git fetch; git checkout {sha}; popd".format(dep=dep, sha=godep_sha))
#   else:
#     print("# UNKNOWN dep %s" % dep)

repos = {}

for godep in godeps['Deps']:
    importpath = godep['ImportPath']
    commit = godep['Rev']
    name = ""

    tokens = importpath.split('/')
    if tokens[0] == "github.com":
        importpath = tokens[0] + "/" + tokens[1] + "/" + tokens[2]

    if tokens[0] == "golang.org" and tokens[1] == "x":
        importpath = tokens[0] + "/" + tokens[1] + "/" + tokens[2]

    if tokens[0] == "gopkg.in":
        importpath = tokens[0] + "/" + tokens[1]

    if tokens[0] == "k8s.io":
        importpath = tokens[0] + "/" + tokens[1]

    if not name:
        name_tokens = []
        for i, t in enumerate(importpath.split('/')):
            if i != 0:
                name_tokens.append(t)
                continue
            for v in reversed(t.split('.')):
                name_tokens.append(v)
        name = '_'.join(name_tokens)

    name = name.replace('.', '_')
    name = name.replace('/', '_')
    name = name.replace('-', '_')

    repo = ('new_go_repository(name="%(name)s", importpath="%(importpath)s", commit="%(commit)s")' %
           {'name': name, 'importpath': importpath, 'commit': commit})
    repos[name] = repo

for name in sorted(repos):
    print repos[name]

# new_go_repository(name="github.com/ugorji/go/codec/codecgen", importpath="github.com/ugorji/go/codec/codecgen", commit="ded73eae5db7e7a0ef6f55aace87a2873c5d2b74")
# new_go_repository(
#     name = "com_github_ugorji_go",
#     importpath = "github.com/ugorji/go",
#     commit = "faddd6128c66c4708f45fdc007f575f75e592a3c",
# )
