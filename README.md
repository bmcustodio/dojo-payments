# dojo-payments

## Testing

In order to run the end-to-end test suite, you may run

```shell
$ make test.e2e
```

In case the API server is not serving requests at `http://localhost:8080`, you must instead run

```shell
$ make test.e2e BASE_URL="http://<host>:<port>"
```

replacing `<host>` and `<port>` with the host and port where the API server can be reached.

## License

Copyright 2019 Bruno Miguel Custodio

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
