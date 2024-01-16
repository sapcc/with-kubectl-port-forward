# with-kubectl-port-forward

As the name implies, this tool executes `kubectl port-forward` for the duration
of the supplied command. For example, Keppel has a PostgreSQL database that can
be reached passwordless on localhost, so the following works to log into the DB:

```
$ with-kubectl-port-forward service/keppel-postgresql 5432:5432 -- psql -U postgres -d keppel -h 127.0.0.1 -p 5432
```

This is the same as:

1. running `kubectl port-forward service/keppel-postgresql 5432:5432` in one shell
2. running `psql -U postgres -d keppel -h 127.0.0.1 -p 5432` in another shell
3. terminating kubectl once psql is done

## Installation

Clone the repo, then run `make install` in it.

## Usage

```
$ with-kubectl-port-forward <port-forward-args>... -- <command>...
```

The exit status will be zero on success, or equivalent to the exit status of
the first failing child otherwise.
