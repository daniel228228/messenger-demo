#!/bin/bash
module_name=$1
major_ver=$2
change_dir=$3

if [[ -z ${module_name} || -z ${major_ver} || -z ${change_dir} ]]; then
  exit 1;
fi

module_name_escaped=$(echo ${module_name} | sed 's/\./\\./g')

if [[ ${major_ver} -eq 1 ]]; then
  if [[ $(go list -m) != ${module_name} ]]; then 
    go mod edit -module "${module_name}"; 
  fi

  find ${change_dir} -name "*.go" -type f -exec sed -i "s#${module_name_escaped}/v\*#${module_name_escaped}#gi" {} \;
else
  if [[ $(go list -m) != ${module_name}/v${major_ver} ]]; then 
    go mod edit -module "${module_name}/v${major_ver}"; 
  fi

  find ${change_dir} -name "*.go" -type f -exec sed -i "s#${module_name_escaped}/v\*#${module_name_escaped}/v${major_ver}#gi" {} \;
fi