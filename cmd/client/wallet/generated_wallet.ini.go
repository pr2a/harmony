package main

const (
	defaultWalletIni = `[default]
bootnode = /ip4/100.26.90.187/tcp/9874/p2p/Qmdfjtk6hPoyrH1zVD9PEH4zfWLo38dP2mDvvKXfh3tnEv
bootnode = /ip4/54.213.43.194/tcp/9874/p2p/QmZJJx6AdaoEkGLrYG4JeLCKeCKDjnFz2wfHNHxAqFSGA9
shards = 4

[default.shard0.rpc]
rpc = l0.t.hmny.io:14555
rpc = s0.t.hmny.io:14555

[default.shard1.rpc]
rpc = l1.t.hmny.io:14555
rpc = s1.t.hmny.io:14555

[default.shard2.rpc]
rpc = l2.t.hmny.io:14555
rpc = s2.t.hmny.io:14555

[default.shard3.rpc]
rpc = l3.t.hmny.io:14555
rpc = s3.t.hmny.io:14555

[local]
bootnode = /ip4/127.0.0.1/tcp/19876/p2p/Qmc1V6W7BwX8Ugb42Ti8RnXF1rY5PF7nnZ6bKBryCgi6cv
shards = 2

[local.shard0.rpc]
rpc = 127.0.0.1:14555
rpc = 127.0.0.1:14557
rpc = 127.0.0.1:14559

[local.shard1.rpc]
rpc = 127.0.0.1:14556
rpc = 127.0.0.1:14558
rpc = 127.0.0.1:14560

[devnet]
bootnode = /ip4/100.26.90.187/tcp/9871/p2p/Qmdfjtk6hPoyrH1zVD9PEH4zfWLo38dP2mDvvKXfh3tnEv
bootnode = /ip4/54.213.43.194/tcp/9871/p2p/QmRVbTpEYup8dSaURZfF6ByrMTSKa4UyUzJhSjahFzRqNj
shards = 4

[devnet.shard0.rpc]
rpc = l0.t1.hmny.io:14555
rpc = s0.t1.hmny.io:14555

[devnet.shard1.rpc]
rpc = l1.t1.hmny.io:14555
rpc = s1.t1.hmny.io:14555

[devnet.shard2.rpc]
rpc = l2.t1.hmny.io:14555
rpc = s2.t1.hmny.io:14555

[devnet.shard3.rpc]
rpc = l3.t1.hmny.io:14555
rpc = s3.t1.hmny.io:14555
`
)
