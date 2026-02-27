# OpenAssembly-Binder
국회 열린데이터광장 에서 제공하는 OpenAPI를 golang 코드로 바인딩해주는 도구입니다.

[국회 열린데이터광장](https://open.assembly.go.kr/portal/openapi/openApiNaListPage.do)에서 제공하는 OpenAPI를 golang 코드로 바인딩해주는 도구입니다.

```
주요한 기능은 동작하지만 아직 미완성입니다.. ;-;
여러가지 기능들과 바인딩 코드의 종류 (python, rust등) 를 향후 추가할 계획입니다!
```

---

## 사용 방법 (빌드)

빌드 전 필요한 프로그램 :
- golang 설치 (1.25.5 이상)
- Task 설치 [설치 방법](https://taskfile.dev/docs/installation#go-modules)

`Taskfile.yml`을 이용하여 

### task build
```
task build
```

### task generate

```
task generate
```

### task list
사용 가능한 api 목록들을 조회합니다.
```
아직 개발중입니다.
```

---
### 출처 및 저작권 표시 (Data Source Attribution)

본 프로젝트는 [국회 열린데이터광장](https://open.assembly.go.kr/)에서 제공하는 OpenAPI를 활용하여 코드를 생성합니다.

* **데이터 제공처**: 대한민국 국회 (National Assembly of the Republic of Korea)
* **이용 조건**: 본 도구가 생성한 코드를 통해 접근하는 데이터는 국회 열린데이터광장의 이용 약관을 따릅니다. 각 API별로 영리적 이용 제한이나 출처 표시 의무가 다를 수 있으니 이용 전 반드시 확인하시기 바랍니다.