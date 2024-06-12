<script>
	const mapbox =
		'pk.eyJ1IjoicmF6b3JzaDRyayIsImEiOiJjbHgyamk4Y3YwZmRiMm1xM2R4bnZ1bDRtIn0.Bt7VcZwdbtfbiDu77xky6g';

	import { onMount } from 'svelte';

	import asta from '$lib/assets/asta.webp';
	import { setModeCurrent } from '@skeletonlabs/skeleton';
	import { AppBar } from '@skeletonlabs/skeleton';
	import { ProgressRadial } from '@skeletonlabs/skeleton';
	import { Icon } from 'svelte-icons-pack';
	import { AiFillGithub } from 'svelte-icons-pack/ai';
	import { TrOutlineBrandGolang, TrOutlineBrandSvelte } from 'svelte-icons-pack/tr';
	import { SiBun, SiGooglecloud, SiMongodb } from 'svelte-icons-pack/si';
	import { Map, Geocoder, Marker, controls } from '@beyonk/svelte-mapbox';

	let tagCloud = [];
	let hMap = [];
	let topValues = [];
	let mapComponent;
	let allRecords = [];
	let selectedTag = '';

	let tagsVisible = false;

	onMount(() => {
		setModeCurrent(false);
		mapComponent.setCenter([-98.43313, 39.50944]);
		mapComponent.setZoom(3);

		fetch('https://indeed-backend-xwdwqsvnia-ue.a.run.app/names')
			.then((res) => res.json())
			.then((res) => (tagCloud = res.names.map((el) => el.replace('_', ' '))))
			.catch((err) => (tagCloud = ['error loading tags']));
	});

	const loadHeatmap = (collection) => {
		selectedTag = collection;

		let url = `https://indeed-backend-xwdwqsvnia-ue.a.run.app/heatmap/${collection.replace(' ', '_')}`;
		fetch(url)
			.then((res) => res.json())
			.then((res) => {
				hMap = res.heatmap;
				minifyHeatmap();
				tagsVisible = true;
			})
			.catch((err) => console.log(err));

		let allUrl = `https://indeed-backend-xwdwqsvnia-ue.a.run.app/records/${collection.replace(' ', '_')}`;
		fetch(allUrl)
			.then((res) => res.json())
			.then((res) => {
				allRecords = res.records.filter((el) => el.longitude && el.latitude);
			})
			.catch((err) => console.log(err));
	};

	const minifyHeatmap = () => {
		if (hMap.length == 0) return [];

		let keys = Object.keys(hMap);
		keys.sort((l, r) => hMap[l] < hMap[r]);

		let objects = keys.slice(0, 30);
		topValues = objects.map((el) => {
			return {
				name: el,
				num: hMap[el]
			};
		});
	};

	const hasIcon = (lang) => {
		const tech = lang.replace(' developer', '').replace(' engineer', '').replace('golang', 'go');
		try {
			const sSheets = document.styleSheets;
			for (let s of sSheets) {
				const rules = s.cssRules;
				for (let r of rules) {
					if (r.cssText.includes(tech)) return true;
				}
			}
			return false;
		} catch (err) {
			return false;
		}
	};

	const getIcon = (lang) => {
		const tech = lang.replace(' developer', '').replace(' engineer', '').replace('golang', 'go');
		return `devicon-${tech}-plain`;
	};
</script>

<div class="grid w-screen md:grid-cols-5">
	<div class="col-span-5 md:col-span-3 md:col-start-2">
		<AppBar gridColumns="grid-cols-3" slotDefault="place-self-center" slotTrail="place-content-end">
			<svelte:fragment slot="lead">
				<img src={asta} alt="logo" width="100wv" />
			</svelte:fragment>
			<div>
				<h1 class="h1 mb-4">Madamada!</h1>
				<p>things you still have to learn</p>
			</div>
			<svelte:fragment slot="trail">
				<a href="https://github.com/RazorSh4rk/job-skill-site" target="_blank">
					<Icon src={AiFillGithub} size="32px" />
				</a>
			</svelte:fragment>
		</AppBar>
	</div>

	<div class="col-span-5 p-2 pt-12 md:col-span-3 md:col-start-2 md:p-0 md:pt-2">
		<p class="h5">Collection of technologies commonly wanted in job listings</p>
		<p class="pt-1">
			Pick a tech from below to see the 30 most requested additional skills for it and the % at
			which they appear in the listings. Every record has a 100+ sample rate. If you don't see
			something but you want to, make a Github issue for it. Clicking on any skill will take you to
			its tutorials.
		</p>

		<div class="grid grid-cols-6 pt-4">
			<p class="col-span-2">Proudly made with</p>
			<div class="col-span-4 grid grid-cols-5 pl-2">
				<Icon src={TrOutlineBrandGolang} size={32} />
				<Icon src={TrOutlineBrandSvelte} size={32} />
				<Icon src={SiBun} size={32} />
				<Icon src={SiMongodb} size={32} />
				<Icon src={SiGooglecloud} size={32} />
			</div>
		</div>
		<hr class="mt-2 !border-t-8 !border-double" />
	</div>

	<div class="col-span-5 pt-12 md:col-span-3 md:col-start-2">
		{#if tagCloud.length == 0}
			<ProgressRadial width="w-8" class="m-auto" />
		{:else}
			<div class="m-auto w-3/4">
				{#each tagCloud as tag}
					<button class="variant-filled badge btn m-2 inline-block p-2" on:click={loadHeatmap(tag)}>
						{#if hasIcon(tag)}
							<i class={getIcon(tag)}></i>
						{/if}
						{tag.replace(' developer', '').replace(' engineer', '')}
					</button>
				{/each}
			</div>
		{/if}

		<hr class="mt-2 !border-t-8 !border-double" />
	</div>

	{#if tagsVisible}
		<div class="col-span-5 overflow-scroll pt-4 md:col-span-3 md:col-start-2">
			{#if topValues.length == 0}
				<ProgressRadial width="w-8" class="m-auto" />
			{:else}
				{#each topValues as skill}
					<div class="relative inline-block">
						<a
							href="https://duckduckgo.com/?t=h_&q={skill.name}+tutorial"
							target="_blank"
							class="variant-filled btn m-2">{skill.name}</a
						>
						<span class="variant-filled badge-icon absolute -right-0 -top-0 p-3"
							>{parseInt((skill.num / Object.keys(hMap).length) * 100)}%</span
						>
					</div>
				{/each}
			{/if}
		</div>
	{/if}
	<div class="col-span-5 overflow-scroll pt-4 md:col-span-3 md:col-start-2">
		<div class="mt-8 h-80 w-full p-2 md:h-96">
			<Map
				accessToken={mapbox}
				bind:this={mapComponent}
				on:recentre={(e) => console.log(e.detail.center.lat, e.detail.center.lng)}
				options={{ scrollZoom: false, doubleClickZoom: false, dragPan: false }}
			>
				{#each allRecords as record}
					<Marker lng={record.longitude} lat={record.latitude} popup={false}></Marker>
				{/each}
			</Map>
		</div>
	</div>
</div>

<!-- <style>
	div {
		border: 1px solid red;
	}
</style> -->
