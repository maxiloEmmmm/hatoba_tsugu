import maxiloVue from 'maxilo-vue'
let utils = maxiloVue.make("utils")

utils.add("btom", (bytes) => {
    if (bytes === 0) return '0 B';
    let k = 1024;
    let sizes = ['B','KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
    let i = Math.floor(Math.log(bytes) / Math.log(k));
    i = Object.is(Infinity, i) || Object.is(-Infinity, i) ? 0 : i
    let num = bytes / Math.pow(k, i);
    return num.toPrecision(3) + ' ' + sizes[i];
})

utils.add("kbid", (str) => {
    //[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*
    return str
        // 小写
        .toLowerCase()
        // 移除其他
        .replace(/[^a-z0-9]/g, "")
        // 去处多余
        .replace(/(^[0-9-]+)/, "")
})

utils.add("kbgitid", (url) => {
    let gitUrl = new URL(url)
    return utils.kbid(gitUrl.pathname)
})

utils.add("kbappid", (id, env) => {
    return `${env}-${id}`
})

utils.add("K8sPathOption", () => {
    return {
        group: "",
        api: "",
        version: "",
        kind: "",
        plural: "",
        ns: ""
    }
})

utils.add("K8sPath", function(option) {
    this.option = Object.assign(utils.K8sPathOption(), option)

    this.multiPath = () => {
        return `/${this.option.group}${this.option.api ? `/${this.option.api}` : ""}/${this.option.version}/namespaces/${this.option.ns}/${this.option.plural}`
    }
    
    this.onePath = (name) =>  {
        return `${this.multiPath()}/${name}`
    }
    
    this.apiVersion = () => {
        return `${this.option.api ? `${this.option.api}/` : ""}${this.option.version}`
    }
})

utils.add("K8sApi", function(path) {
    return {
        list: (labelSelector = "") => {
            return maxiloVue.app.$kb.get(path.multiPath(), {params: {
                labelSelector: labelSelector
            }})
        },
        get: (name, labelSelector = "") => {
            return maxiloVue.app.$kb.get(path.onePath(name), {params: {
                labelSelector: labelSelector
            }})
        },
        update: (data) => {
            return maxiloVue.app.$kb.put(path.onePath(data.metadata.name), data)
        },
        create: (data) => {
            return maxiloVue.app.$kb.post(path.multiPath(), data)
        },
        delete: (name) => {
            return maxiloVue.app.$kb.delete(path.onePath(name))
        },
        deleteLabel: (labelSelector = "") => {
            return maxiloVue.app.$kb.delete(path.multiPath(), {params: {
                labelSelector: labelSelector
            }})
        },
        fullUpdateOrCreate: async (data) => {
            try {
                let response = await maxiloVue.app.$kb.get(path.onePath(data.metadata.name))
                data.metadata.resourceVersion = response.data.metadata.resourceVersion
                await maxiloVue.app.$kb.put(path.onePath(data.metadata.name), data)
            }catch(e) {
                if(e.response.data.code == 404) {
                    try {
                        await maxiloVue.app.$kb.post(path.multiPath(), data)
                    }catch(e) {
                        throw e
                    }
                }else {
                    throw e
                }
            }
        },
        url: {
            one: path.onePath,
            multi: path.multiPath()
        },
        path,
    }
})