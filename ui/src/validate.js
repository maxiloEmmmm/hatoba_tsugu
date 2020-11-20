import maxiloVue from "maxilo-vue"
import { required } from 'vee-validate/dist/rules';

maxiloVue.register({
    register: function(app){
        const validator = app.make("validator")
        validator.addRule("required", {
            ...required,
            message: "必填不可为空!"
        })
        
        validator.addRule("abc", {
            validate(value){
                return /^[a-zA-Z0-9-_]+$/.test(value)
            },
            message: "只可以为大小写字母数字和-_"
        })
    },
    boot(){}
})