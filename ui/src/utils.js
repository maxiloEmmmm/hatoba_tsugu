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
        .replace(/[^a-z0-9-]/g, "-")
        // 去处多余
        .replace(/(^[0-9-]+|[-]+$|[-]{2,})/, "")
})