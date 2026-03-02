# OpenAssembly-Binder

[국회 열린데이터광장](https://open.assembly.go.kr/portal/openapi/openApiNaListPage.do)에서 제공하는 OpenAPI를 Go 코드로 자동 바인딩해주는 도구입니다.

> **상태**: 핵심 기능은 동작하며, 추가 기능과 다른 언어 지원(Python, Rust 등)을 계획 중입니다.

---

## 요구사항

- **Go** 1.25.5 이상
- **Task** - [설치 가이드](https://taskfile.dev/docs/installation#go-modules)

## 사용 방법

### 빌드
```bash
task build
```

### 코드 생성
```bash
task generate
```
OpenAPI 스펙을 읽어 Go 바인딩 코드를 자동 생성합니다.

### API 목록 조회
```bash
task list
```
사용 가능한 API 목록을 조회합니다. (개발 진행 중)

---

## 라이선스 및 출처

본 프로젝트는 [국회 열린데이터광장](https://open.assembly.go.kr/)에서 제공하는 OpenAPI를 활용합니다.

**데이터 제공처**: 대한민국 국회 (National Assembly of the Republic of Korea)

**이용 약관**: 본 도구가 생성한 코드를 통해 접근하는 데이터는 국회 열린데이터광장의 이용 약관을 따릅니다. API별로 영리적 이용 제한이나 출처 표시 의무가 다를 수 있으므로 이용 전 반드시 확인하시기 바랍니다.