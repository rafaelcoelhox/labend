# 🚀 GraphQL - Exemplos de Uso

Este documento contém exemplos de queries e mutations GraphQL para testar a aplicação LabEnd.

## 🌐 Endpoints

- **GraphQL API**: `http://localhost:8080/graphql`
- **GraphQL Playground**: `http://localhost:8080/graphql` (GET)

## 👥 Usuários

### Queries

```graphql
# Buscar usuário por ID
query GetUser {
  user(id: "1") {
    id
    name
    email
    totalXP
    createdAt
    updatedAt
  }
}

# Listar usuários com paginação
query GetUsers {
  users(limit: 5, offset: 0) {
    id
    name
    email
    totalXP
    createdAt
  }
}

# Histórico de XP do usuário
query GetUserXPHistory {
  userXPHistory(userID: "1") {
    id
    userID
    sourceType
    sourceID
    amount
    createdAt
  }
}
```

### Mutations

```graphql
# Criar usuário
mutation CreateUser {
  createUser(input: {
    name: "João Silva"
    email: "joao@exemplo.com"
  }) {
    id
    name
    email
    totalXP
    createdAt
  }
}

# Atualizar usuário
mutation UpdateUser {
  updateUser(id: "1", input: {
    name: "João Santos"
    email: "joao.santos@exemplo.com"
  }) {
    id
    name
    email
    totalXP
    updatedAt
  }
}

# Deletar usuário
mutation DeleteUser {
  deleteUser(id: "1")
}
```

## 🏆 Challenges

### Queries

```graphql
# Buscar challenge por ID
query GetChallenge {
  challenge(id: "1") {
    id
    title
    description
    xpReward
    status
    createdAt
    updatedAt
  }
}

# Listar challenges
query GetChallenges {
  challenges(limit: 10, offset: 0) {
    id
    title
    description
    xpReward
    status
    createdAt
  }
}

# Submissões de um challenge
query GetChallengeSubmissions {
  challengeSubmissions(challengeID: "1") {
    id
    challengeID
    userID
    proofURL
    status
    createdAt
  }
}

# Votos de uma submissão
query GetChallengeVotes {
  challengeVotes(submissionID: "1") {
    id
    submissionID
    userID
    approved
    timeCheck
    isValid
    createdAt
  }
}
```

### Mutations

```graphql
# Criar challenge
mutation CreateChallenge {
  createChallenge(input: {
    title: "Aprender GraphQL"
    description: "Complete um tutorial de GraphQL e crie uma API"
    xpReward: 150
  }) {
    id
    title
    description
    xpReward
    status
    createdAt
  }
}

# Submeter challenge
mutation SubmitChallenge {
  submitChallenge(input: {
    challengeID: "1"
    proofURL: "https://github.com/user/graphql-project"
  }) {
    id
    challengeID
    userID
    proofURL
    status
    createdAt
  }
}

# Votar em submissão
mutation VoteChallenge {
  voteChallenge(input: {
    submissionID: "1"
    approved: true
    timeCheck: 2500
  }) {
    id
    submissionID
    userID
    approved
    timeCheck
    isValid
    createdAt
  }
}
```

## 🔄 Queries Combinadas

```graphql
# Query complexa - Usuário com challenges e XP
query GetUserComplete {
  user(id: "1") {
    id
    name
    email
    totalXP
    createdAt
  }
  
  challenges(limit: 5) {
    id
    title
    xpReward
    status
  }
  
  userXPHistory(userID: "1") {
    sourceType
    sourceID
    amount
    createdAt
  }
}
```

## 🧪 Testes de Performance

```graphql
# Query otimizada - usando JOIN para eliminar N+1
query GetUsersOptimized {
  users(limit: 100, offset: 0) {
    id
    name
    email
    totalXP  # <- Calculado via JOIN, não N+1 queries
    createdAt
  }
}
```

## 📊 Queries de Analytics

```graphql
# Estatísticas gerais
query GetStats {
  users(limit: 1) {
    id
  }
  
  challenges(limit: 1) {
    id
  }
  
  # Top usuários por XP
  users(limit: 10, offset: 0) {
    name
    totalXP
  }
}
```

## 🔍 Filtros e Busca

```graphql
# Buscar challenges ativos
query GetActiveChallenges {
  challenges(limit: 20) {
    id
    title
    xpReward
    status
    createdAt
  }
}
```

## 🚀 Como Testar

### 1. Via GraphQL Playground

1. Acesse `http://localhost:8080/graphql`
2. Cole uma query/mutation
3. Clique em "Play"

### 2. Via cURL

```bash
# Query via cURL
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "query GetUsers { users(limit: 5) { id name email totalXP } }"
  }'

# Mutation via cURL
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation CreateUser($input: CreateUserInput!) { createUser(input: $input) { id name email } }",
    "variables": {
      "input": {
        "name": "Test User",
        "email": "test@example.com"
      }
    }
  }'
```

### 3. Via Cliente JavaScript

```javascript
// Usando fetch
const query = `
  query GetUsers {
    users(limit: 5) {
      id
      name
      email
      totalXP
    }
  }
`;

fetch('http://localhost:8080/graphql', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({ query }),
})
.then(response => response.json())
.then(data => console.log(data));
```

## 📝 Notas Importantes

1. **IDs são strings** no GraphQL (convertidos internamente)
2. **XP é calculado** via JOIN otimizada (sem N+1)
3. **Campos opcionais** podem ser omitidos
4. **Validações** são aplicadas nos inputs
5. **Erros** são retornados no formato GraphQL padrão

## 🎯 Casos de Uso Avançados

### Subscription (Futuro)

```graphql
# Para implementar no futuro
subscription OnChallengeSubmitted {
  challengeSubmitted {
    id
    challengeID
    userID
    proofURL
  }
}
```

### Fragments

```graphql
fragment UserInfo on User {
  id
  name
  email
  totalXP
}

query GetUserWithFragment {
  user(id: "1") {
    ...UserInfo
    createdAt
  }
}
```

---

**Resultado**: Agora você tem exemplos completos para testar toda a funcionalidade GraphQL da aplicação! 