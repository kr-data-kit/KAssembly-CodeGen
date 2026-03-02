# KAssembly-CodeGen

[국회 열린데이터광장](https://open.assembly.go.kr/portal/openapi/openApiNaListPage.do)에서 제공하는 OpenAPI의 클라이언트 코드를 생성합니다.

**지원 언어**: Go, Python

> **상태**: 핵심 기능은 동작하며, 추가 기능을 개발 중입니다. (beta)

---

## 요구사항

- **Go** 1.25.5 이상

## 설치

```bash
# 저장소 클론
git clone https://github.com/kr-data-kit/KAssembly-CodeGen
cd KAssembly-CodeGen

# 의존성 설치 및 빌드
go mod download
go build -o ./build/kassemblycodegen.exe .
```

## 사용 방법
### 코드 생성

**Go 클라이언트 생성:**
```bash
./build/kassemblycodegen generate -m go
```

**Python 클라이언트 생성:**
```bash
./build/kassemblycodegen generate -m python
```

**옵션 지정:**
```bash
./build/kassemblycodegen generate -m go \
  --package myauth \
  --client MyClient \
  --output ./generated \
  --create-dir
```

| 옵션 | 기본값 | 설명 |
|------|--------|------|
| `-m, --language` | - | 생성할 언어 (go, python) **필수** |
| `--package` | openassemblyclient | 패키지/모듈 이름 |
| `--client` | OpenAssemblyClient | 클라이언트 클래스/구조체 이름 (Go만) |
| `--output` | ./out | 출력 디렉토리 |
| `--create-dir` | false | 출력 디렉토리가 없으면 생성 |

### API 목록 조회
```bash
./build/kassemblycodegen list --method simple
```

---

## 예제

### Go 클라이언트 생성 및 사용
```bash
# 빌드
go build -o ./build/kassemblycodegen.exe .

# 생성
./build/kassemblycodegen generate -m go --package openassembly --output ./generated --create-dir

# 생성된 코드는 ./generated 디렉토리에 위치
```

```go
package main

import (
  "context"
  "log"

  "your-module/generated/openassembly"
)

func main() {
  client := openassembly.NewOpenAssemblyClient("YOUR_API_KEY")

  // 예: 생성된 Fetch{Service} 메서드 사용
  resp, err := client.FetchAllBill(context.Background(), &openassembly.AllBillParams{})
  if err != nil {
    log.Fatal(err)
  }

  log.Printf("status=%s rows=%d", resp.Status, len(resp.Data))
}
```

> 참고: 서비스마다 메서드 이름은 `Fetch{Service}` 형태로 생성됩니다.

### Python 클라이언트 생성 및 사용
```bash
# 빌드
go build -o ./build/kassemblycodegen.exe .

# 생성
./build/kassemblycodegen generate -m python --package openassemblyclient --output ./generated --create-dir

# 의존성 설치
cd generated
uv sync

# 사용
from openassemblyclient import Client
client = Client(api_key="YOUR_API_KEY")
```

---

## 라이선스 및 출처

본 프로젝트는 [국회 열린데이터광장](https://open.assembly.go.kr/)에서 제공하는 OpenAPI를 활용합니다.

**데이터 제공처**: 대한민국 국회 (National Assembly of the Republic of Korea)

**이용 약관**: 본 도구가 생성한 코드를 통해 접근하는 데이터는 국회 열린데이터광장의 이용 약관을 따릅니다. API별로 영리적 이용 제한이나 출처 표시 의무가 다를 수 있으므로 이용 전 반드시 확인하시기 바랍니다.