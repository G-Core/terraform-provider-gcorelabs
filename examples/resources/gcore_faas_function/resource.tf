provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_faas_function" "func" {
        project_id = 1
        region_id = 1
        name = "testf"
        namespace = "ns4test"
        description = "function description"
        envs = {
                BIG = "EXAMPLE2"
        }
        runtime = "go1.16.6"
        code_text = <<EOF
package kubeless

import (
        "github.com/kubeless/kubeless/pkg/functions"
)

func Run(evt functions.Event, ctx functions.Context) (string, error) {
        return "Hello World!!", nil
}
EOF
        timeout = 5
        flavor = "80mCPU-128MB"
        main_method = "main"
        min_instances = 1
        max_instances = 2
}
