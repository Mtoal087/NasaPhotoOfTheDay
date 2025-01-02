(async () => {
  await loadLinksPreset(tsParticles);

  await tsParticles.load({
    id: "tsparticles",
    options: {
      preset: "links",
    },
  });
})();