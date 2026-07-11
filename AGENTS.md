# Agent instructions

## Code review (CodeRabbit / PR comments)

When addressing review comments on a pull request (especially CodeRabbit):

1. **Fix** the issue in code (and tests when needed).
2. **Reply** on the thread with a short note of what changed (commit SHA when useful).
3. **Resolve** the GitHub review thread after it is addressed — do not leave it open once the fix is pushed and explained.

### Resolving threads

Use the GraphQL `resolveReviewThread` mutation with the thread’s node id (`PRRT_...`), for example:

```bash
gh api graphql -f query='
mutation($id: ID!) {
  resolveReviewThread(input: { threadId: $id }) {
    thread { isResolved }
  }
}' -F id='PRRT_...'
```

List unresolved threads for a PR:

```bash
gh api graphql -f query='
query($n: Int!) {
  repository(owner: "OWNER", name: "REPO") {
    pullRequest(number: $n) {
      reviewThreads(first: 100) {
        nodes {
          id
          isResolved
          comments(first: 1) { nodes { author { login } path } }
        }
      }
    }
  }
}' -F n=PR_NUMBER
```

Only leave a thread unresolved if the user explicitly wants discussion kept open, or the finding is deferred with a clear reply explaining why.
