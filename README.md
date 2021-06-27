# arc

Arc is my knock off version of of the Phabricator arcanist tooling for a similar experience when using  branch based workflows.

### flow

The flow command shows you a graph of the branches you are working on locally. If
the branch is tracking another local branch or has other local branches tracking it
flow will display it. Branches that are tracking remotes only will not be shown.

In short it will take the output of `git branch -vvv` and turn it from this

```
  gone                                                              25503de [main: behind 1] Properly guard against npe
* main                                                              3f126e5 Cleanup issue for branch with missing merge point
  mani                                                              5232f9f [maini: gone] Cleanup status bar and make it smooth
  master                                                            539c188 wowowowowo
  mschuett/arc-flow-init-cmd                                        5232f9f [main: behind 8] Cleanup status bar and make it smooth
  mschuett/better-branch                                            5232f9f [main: behind 8] Cleanup status bar and make it smooth
  mschuett/made-from-flow                                           5232f9f [main: behind 8] Cleanup status bar and make it smooth
  mschuett/off-master                                               539c188 [master] wowowowowo
  mschuett/off-master-2                                             539c188 [mschuett/off-master] wowowowowo
  mschuett/off-test-2                                               7df606e [mschuett/test-2] Test commit
  mschuett/test-2                                                   7df606e [main: ahead 1, behind 8] Test commit
  mschuett/test-branch                                              5232f9f [mschuett/arc-flow-init-cmd] Cleanup status bar and make it smooth
  mschuett/testing                                                  5232f9f [main: behind 8] Cleanup status bar and make it smooth
  mschuett/testing-deeper-1                                         7df606e [mschuett/test-2] Test commit
  mschuett/testing-deeper-2                                         7df606e [mschuett/testing-deeper-1] Test commit
  mschuett/testing-deeper-3                                         7df606e [mschuett/testing-deeper-2] Test commit
  mschuett/testtttttttt                                             5232f9f [main: behind 8] Cleanup status bar and make it smooth
  mschuett/trash-test                                               c9d9853 [master: ahead 1, behind 1] lol
  mschuett/whaaa                                                    5232f9f [main: behind 8] Cleanup status bar and make it smooth
  mschuett/wow-this-branch-name-is-really-long-it-doesnt-need-to-be 5232f9f [main: behind 8] Cleanup status bar and make it smooth
```

into this

```
master
 ├ mschuett/trash-test                                                 c9d98532 1:1 lol
 └ mschuett/off-master                                                 539c188a 0:0 wowowowowo
  └ mschuett/off-master-2                                              539c188a 0:0 wowowowowo
main ๏
 ├ mschuett/wow-this-branch-name-is-really-long-it-doesnt-need-to-be   5232f9fd 8:0 Cleanup status bar and make it smooth
 ├ mschuett/whaaa                                                      5232f9fd 8:0 Cleanup status bar and make it smooth
 ├ mschuett/testtttttttt                                               5232f9fd 8:0 Cleanup status bar and make it smooth
 ├ mschuett/testing                                                    5232f9fd 8:0 Cleanup status bar and make it smooth
 ├ mschuett/test-2                                                     7df606e2 8:1 Test commit
 │├ mschuett/testing-deeper-1                                          7df606e2 0:0 Test commit
 ││└ mschuett/testing-deeper-2                                         7df606e2 0:0 Test commit
 ││ └ mschuett/testing-deeper-3                                        7df606e2 0:0 Test commit
 │└ mschuett/off-test-2                                                7df606e2 0:0 Test commit
 ├ mschuett/made-from-flow                                             5232f9fd 8:0 Cleanup status bar and make it smooth
 ├ mschuett/better-branch                                              5232f9fd 8:0 Cleanup status bar and make it smooth
 ├ mschuett/arc-flow-init-cmd                                          5232f9fd 8:0 Cleanup status bar and make it smooth
 │└ mschuett/test-branch                                               5232f9fd 0:0 Cleanup status bar and make it smooth
 └ gone                                                                25503de1 1:0 Properly guard against npe
```
