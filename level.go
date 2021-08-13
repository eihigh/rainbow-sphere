package main

import "github.com/eihigh/rainbow-sphere/model"

var (
	configs = []*model.Config{
		{ // stage 1
			HP:       5,
			SphereHP: 10,
			MinSpeed: 3,
			AmpSpeed: 1,
			MinScale: 1,
			AmpScale: 0,
		},
		{ // stage 2
			HP:       5,
			SphereHP: 20,
			MinSpeed: 4,
			AmpSpeed: 2,
			MinScale: 1,
			AmpScale: 0,
		},
		{ // stage 3
			HP:       5,
			SphereHP: 25,
			MinSpeed: 5,
			AmpSpeed: 5,
			MinScale: 1,
			AmpScale: 0,
		},
		{ // stage 4
			HP:       5,
			SphereHP: 50,
			MinSpeed: 5,
			AmpSpeed: 2,
			MinScale: 3.5,
			AmpScale: 0,
		},
		{ // stage 5
			HP:       5,
			SphereHP: 30,
			MinSpeed: 3,
			AmpSpeed: 8,
			MinScale: 1,
			AmpScale: 2.5,
		},
		// { // stage 6
		// 	HP:       5,
		// 	SphereHP: 50,
		// 	MinSpeed: 7,
		// 	AmpSpeed: 2,
		// 	MinScale: 3,
		// 	AmpScale: 0.5,
		// },
		{ // stage 7
			HP:       5,
			SphereHP: 25,
			MinSpeed: 5,
			AmpSpeed: 10,
			MinScale: 1,
			AmpScale: 0,
		},
	}

	// 一定数クリア後はここからランダムに
	endlessConfigs = []*model.Config{
		{ // stage 4
			SphereHP: 50,
			MinSpeed: 6,
			AmpSpeed: 5,
			MinScale: 3.8,
			AmpScale: 0,
		},
		{ // stage 5
			SphereHP: 35,
			MinSpeed: 7,
			AmpSpeed: 8,
			MinScale: 1,
			AmpScale: 2.5,
		},
		{ // stage 6
			SphereHP: 50,
			MinSpeed: 9,
			AmpSpeed: 3,
			MinScale: 3,
			AmpScale: 1,
		},
		{ // stage 7
			SphereHP: 25,
			MinSpeed: 5,
			AmpSpeed: 11,
			MinScale: 1,
			AmpScale: 0,
		},
	}
)
