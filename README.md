# arc

Arc is my knock off version of of the Phabricator arcanist tooling for a similar experience when using  branch based workflows.

### flow

The flow command shows you a graph of the branches you are working on locally. If
the branch is tracking another local branch or has other local branches tracking it
flow will display it. Branches that are tracking remotes only will not be shown.

In short it will take the output of `git branch -vvv` and turn it from this

```
  main                       3815e64 Outline of arcinist cli
  master                     3815e64 Outline of arcinist cli
  mschuett/arc-flow-init-cmd 3815e64 [main] Outline of arcinist cli
  mschuett/off-master        3815e64 [master] Outline of arcinist cli
* mschuett/off-master-2      3815e64 [mschuett/off-master] Outline of arcinist cli
  mschuett/test-branch       3815e64 [mschuett/arc-flow-init-cmd] Outline of arcinist cli
  mschuett/testing           3815e64 [main] Outline of arcinist cli
```

into this

```
main
  ├─ mschuett/testing
  └─ mschuett/arc-flow-init-cmd
     └─ mschuett/test-branch

master
  mschuett/off-master
    mschuett/off-master-2
```
