<h1 align="center">Arcanist CLI</h1>

<p align="center">Git workflow tool based on the Phabricator workflow</p>

<br/>

Arcanist is an local git workflow management tool that is based around local branch management. This allows you to easily keep ideas separate and swich between branches pushing the changes up to other branches with `cascade`. Using this tool should help to keep PRs/Diffs/Whatever to smaller sizes that are more easily reviewable.

## flow [branch] [upstream]

Outputs a graph of all your current branches that are tracking local branches to the CLI. When you pass a single argument to flow it will create a new branch and set it's upstream to the current branch. If the branch exists it will just check out that branch for you. If you pass two arguments to arc flow it will create a new branch and set it's upstream to the second argument instead of to the branch you are currently on.

<img width="571" alt="Screen Shot 2021-07-01 at 6 56 21 PM" src="https://user-images.githubusercontent.com/2604634/124201891-cf261700-da9e-11eb-826d-4a35795897c1.png">

## cascade

Rebase all branches that are downstreams of the branch you are currently on. If the rebase has conflicts it will abort the rebase and keep going letting the user know at the end which branches require manual intervention.

<img width="573" alt="Screen Shot 2021-07-01 at 6 57 42 PM" src="https://user-images.githubusercontent.com/2604634/124201808-9423e380-da9e-11eb-8165-8078e47e48e3.png">

## prune

Delete all branches where all commits in the branch are already contained in the parent. It will additionally try to re-parent all branches so the branch doesn't fall out of the `arc flow` output. This may not be possible is all branches upstream of a branch are deleted.

<img width="560" alt="Screen Shot 2021-07-01 at 6 58 58 PM" src="https://user-images.githubusercontent.com/2604634/124201775-78b8d880-da9e-11eb-9097-0a3ee25db190.png">

## graft

Pull a branch from a given remote into your working tree.

## diff / commit

coming soon.
