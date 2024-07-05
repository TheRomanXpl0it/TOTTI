const flagCtx = document.getElementById('SubmittedFlags').getContext('2d');
const exploitCtx = document.getElementById('Exploits').getContext('2d');
const teamCtx = document.getElementById('Teams').getContext('2d');
const roundCtx = document.getElementById('Rounds').getContext('2d');
const exploits = document.getElementById('exploit');
const manualFlag = document.getElementById('flagInput');
window.startingTick = 0;
window.lastTick = 0;
document.cookie = "startTick=";
document.cookie = "exploit=";

function createOrUpdateCharts(data) {
	console.log(data);

	window.lastTick = Object.keys(data.rounds).length;
	/*let cookie = document.cookie.split(';').find(c => c.trim().startsWith('startTick=')).split('=')[1];
	console.log(cookie);
	if (cookie[0] === 'l') {
		window.startingTick =  window.lastTick - cookie.substring(1);
	} else {
		window.startingTick = cookie;
	}
	console.log(window.startingTick, window.lastTick);*/

	const flagData = data.flags;
	const exploitData = data.exploits;
	const teamData = data.teams;
	const roundData = data.rounds;



	if (!window.rounds) {
		window.rounds = new Chart(roundCtx, {
			type: 'line',
			data: {
				labels: [],
				datasets: [{
					label: 'Accepted',
					borderColor: '#b372ff',
					//cubicInterpolationMode: 'monotone',
					//tension: 0.4,
					data: []
				}, {
					type: 'bar',
					stack: 'combined',
					label: 'Accepted',
					backgroundColor: '#93c54b',
					borderColor: '#597e26',
					borderWidth: 1,
					barPercentage: 1.0,
					categoryPercentage: 0.8,
					data: []
				}, {
					type: 'bar',
					stack: 'combined',
					label: 'Error',
					backgroundColor: '#d9534f',
					borderColor: '#9e3431',
					borderWidth: 1,
					barPercentage: 1.0,
					categoryPercentage: 0.8,
					data: []
				}, {
					type: 'bar',
					stack: 'combined',
					label: 'Expired',
					backgroundColor: '#f47c3c',
					borderColor: '#3f1804',
					borderWidth: 1,
					barPercentage: 1.0,
					categoryPercentage: 0.8,
					data: []
				}, {
					type: 'bar',
					stack: 'combined',
					label: 'Queued',
					backgroundColor: '#29abe0',
					borderColor: '#0074a2',
					borderWidth: 1,
					barPercentage: 1.0,
					categoryPercentage: 0.8,
					data: []
				}]
			},
			options: {
				responsive: true,
				plugins: {
					legend: {
						display: false,
						position: 'top',
						labels: {
							boxWidth: 20
						}
					}
				},
				scales: {
					y: {
						stacked: true,
					},
					x: {
						stacked: true,
					},
				},
			}
		});
	}

	{
		for (let i = window.startingTick; i < Object.keys(roundData).length; i++) {
			let item = roundData[i];

			if (item === undefined) {
				item = {accepted: 0, error: 0, expired: 0, queued: 0};
			}

			if (window.rounds.data.labels.includes(i)) {

				window.rounds.data.datasets[0].data[i - window.startingTick].y = item.accepted + item.error + item.expired + item.queued;
				window.rounds.data.datasets[1].data[i - window.startingTick].y = item.accepted;
				window.rounds.data.datasets[2].data[i - window.startingTick].y = item.error;
				window.rounds.data.datasets[3].data[i - window.startingTick].y = item.expired;
				window.rounds.data.datasets[4].data[i - window.startingTick].y = item.queued;

			} else {

				window.rounds.data.labels.splice(i - window.startingTick, 0, i);
				window.rounds.data.datasets[0].data.splice(i - window.startingTick, 0, {x: i, y: item.accepted + item.error + item.expired + item.queued});
				window.rounds.data.datasets[1].data.splice(i - window.startingTick, 0, {x: i, y: item.accepted});
				window.rounds.data.datasets[2].data.splice(i - window.startingTick, 0, {x: i, y: item.error});
				window.rounds.data.datasets[3].data.splice(i - window.startingTick, 0, {x: i, y: item.expired});
				window.rounds.data.datasets[4].data.splice(i - window.startingTick, 0, {x: i, y: item.queued});
			}
		}

		window.rounds.update();
	}



	if (window.teams) {

		let i = 0;
		for (const key of Object.keys(teamData).sort((a, b) => a.length - b.length)) {
			let item = teamData[key];
			window.teams.data.datasets[0].data[i].y = item.accepted;
			window.teams.data.datasets[1].data[i].y = item.error;
			window.teams.data.datasets[2].data[i].y = item.expired;
			window.teams.data.datasets[3].data[i].y = item.queued;
			i++;
		};

		window.teams.update();
	} else {

		window.teams = new Chart(teamCtx, {
			type: 'bar',
			data: {
				labels: [],
				datasets: [{
					label: 'Accepted',
					backgroundColor: '#93c54b',
					borderColor: '#597e26',
					borderWidth: 1,
					barPercentage: 1.0,
					categoryPercentage: 0.8,
					data: []
				}, {
					label: 'Error',
					backgroundColor: '#d9534f',
					borderColor: '#9e3431',
					borderWidth: 1,
					barPercentage: 1.0,
					categoryPercentage: 0.8,
					data: []
				}, {
					label: 'Expired',
					backgroundColor: '#f47c3c',
					borderColor: '#3f1804',
					borderWidth: 1,
					barPercentage: 1.0,
					categoryPercentage: 0.8,
					data: []
				}, {
					label: 'Queued',
					backgroundColor: '#29abe0',
					borderColor: '#0074a2',
					borderWidth: 1,
					barPercentage: 1.0,
					categoryPercentage: 0.8,
					data: []
				}]
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				plugins: {
					legend: {
						display: false,
						position: 'top',
						labels: {
							boxWidth: 20
						}
					}
				},
				scales: {
					y: {
						stacked: true,
					},
					x: {
						stacked: true,
					},
				},
			}
		});

		let i = 0;
		for (const key of Object.keys(teamData).sort((a, b) => a.length - b.length)) {
			let item = teamData[key];
			window.teams.data.labels.splice(i, 0, key);
			window.teams.data.datasets[0].data.splice(i, 0, {x: key, y: item.accepted});
			window.teams.data.datasets[1].data.splice(i, 0, {x: key, y: item.error});
			window.teams.data.datasets[2].data.splice(i, 0, {x: key, y: item.expired});
			window.teams.data.datasets[3].data.splice(i, 0, {x: key, y: item.queued});
			i++;
		};
		window.teams.update();
	}



	if (!window.exploits) {
		window.exploits = new Chart(exploitCtx, {
			type: 'bar',
			data: {
				labels: [],
				datasets: [{
					label: 'Accepted',
					backgroundColor: '#93c54b',
					barPercentage: 0.7,
					categoryPercentage: 0.8,
					data: []
				}, {
					label: 'Queued',
					backgroundColor: '#29abe0',
					barPercentage: 0.7,
					categoryPercentage: 0.8,
					data: []
				}, {
					label: 'Expired',
					backgroundColor: '#f47c3c',
					barPercentage: 0.7,
					categoryPercentage: 0.8,
					data: []
				}, {
					label: 'Error',
					backgroundColor: '#d9534f',
					barPercentage: 0.7,
					categoryPercentage: 0.8,
					data: []
				}]
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				plugins: {
					legend: {
						display: false,
						position: 'top',
						labels: {
							boxWidth: 20
						}
					}
				}
			}
		});
	}

	{
		let i = 0;
		for (const [key, item] of Object.entries(exploitData)) {
			if (window.exploits.data.labels.includes(key)) {
				
				let idxAccepted = window.exploits.data.datasets[0].data.findIndex(obj => {
					return obj.x === key;
				});
				let idxQueued = window.exploits.data.datasets[1].data.findIndex(obj => {
					return obj.x === key;
				});
				let idxExpired = window.exploits.data.datasets[2].data.findIndex(obj => {
					return obj.x === key;
				});
				let idxError = window.exploits.data.datasets[3].data.findIndex(obj => {
					return obj.x === key;
				});

				window.exploits.data.datasets[0].data[idxAccepted].y = item.accepted;
				window.exploits.data.datasets[1].data[idxQueued].y = item.queued;
				window.exploits.data.datasets[2].data[idxExpired].y = item.expired;
				window.exploits.data.datasets[3].data[idxError].y = item.error;

			} else {

				window.exploits.data.labels.splice(i, 0, key);
				let acceptedObj = {x: key, y: item.accepted}
				let queuedObj = {x: key, y: item.queued}
				let expiredObj = {x: key, y: item.expired}
				let errorObj = {x: key, y: item.error}
				window.exploits.data.datasets[0].data.splice(i, 0, acceptedObj);
				window.exploits.data.datasets[1].data.splice(i, 0, queuedObj);
				window.exploits.data.datasets[2].data.splice(i, 0, expiredObj);
				window.exploits.data.datasets[3].data.splice(i, 0, errorObj);

				var opt = document.createElement("option");
				opt.value = key;
				opt.innerHTML = key;
				exploits.appendChild(opt);
			}
			i++;
		}

		window.exploits.update();
	}



	if (window.submittedFlags) {
		window.submittedFlags.data.datasets[0].data = [flagData.accepted, flagData.queued, flagData.expired, flagData.error];
		window.submittedFlags.update();
	} else {
		window.submittedFlags = new Chart(flagCtx, {
			type: 'doughnut',
			data: {
				labels: ['Accepted', 'Queued', 'Expired', 'Error'],
				datasets: [{
					label: 'Dataset',
					data: [flagData.accepted, flagData.queued, flagData.expired, flagData.error],
					borderWidth: 1,
					backgroundColor: ['#93c54b', '#29abe0', '#f47c3c', '#d9534f'],
				}]
			},
			options: {
				responsive: true,
				plugins: {
					legend: {
						position: 'bottom',
						labels: {
							color: 'white',
						},
					},
				},
			}
		});
	}
}

document.addEventListener('htmx:afterOnLoad', (event) => {
	if (event.detail.pathInfo.path === '/data') {
		createOrUpdateCharts(JSON.parse(event.detail.xhr.responseText));
	}
});

function changeStart(value) {
	if (value[0] === 'l') {
		window.startingTick =  window.lastTick - value.substring(1);
	} else {
		window.startingTick = value;
	}
	document.cookie = "startTick="+value;
	window.rounds.destroy();
	window.rounds = null;
	fetch('/data')
		.then(response => response.json())
		.then(data => createOrUpdateCharts(data));
}

function changeExploit(value) {
	document.cookie = "exploit="+value;
	fetch('/data')
		.then(response => response.json())
		.then(data => createOrUpdateCharts(data));
}

function submitFlag() {
	const flagValue = manualFlag.value;
	fetch('/manual?flag='+flagValue, {method: 'POST'})
}
