# 중요 TODO

- 패키지 이름이 가능한지 체크하는 절차 필요

- 테스트 코드를 만들어야 함

## CLI

- `ctx` config 추가

- golang 생성시, go.mod 생성 여부 설정

- 특정 endpoint 만 생성하거나 생성 하지 않게 하는 기능 추가 : 성능 개선하기 (`include`제약시에도 모든 api를 가져옴)

- `generate` 진행과정 표시

---

# Future

- example 자동 생성

- service 구조체가 값이 빠져도 gen이 가능하게 설계를 바꾸기

- `common.go` 손보기 (`xml` 제거) : xml로 바꾸는 것은 현재 불필요하지만, 없애는 게 급선무는 아님

- bulk update (xlsx) : 외부 패키지 필요 + json에 비해 어려움

- 탬플릿에서 ast로 변환하기 : 보안 관련, 급선무가 아님

- 더 나은 명명 (en-name) + 카테고리화 : 초반에 추가했으나, 중요한 API는 이름이 잘 구분되어있고, en 생성 퀄리티가 낮아서 나중에 다시 하기로