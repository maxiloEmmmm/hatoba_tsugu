import maxiloVue from 'maxilo-vue'
let router = maxiloVue.make("router").getRoute()

router.add("/", () => import("./pages/test.vue"))
router.group("/hatoba_tsugu_deploy", "", function(r){
    r.add("project", () => import("./pages/hatoba_tsugu/deploy/project.vue"))
    r.add("project/:id", () => import("./pages/hatoba_tsugu/deploy/space.vue"))
})