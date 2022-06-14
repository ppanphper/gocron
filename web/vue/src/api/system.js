import httpClient from '../utils/httpClient'

export default {
    loginLogList(query, callback) {
        httpClient.get('/system/login-log', query, callback)
    },
    updateSystemSetting(data, callback) {
        httpClient.post('/system/setting', data, callback)
    },
    getSystemSetting(callback) {
        httpClient.get('/system/setting', {}, callback)
    },
    resetSystemSetting(app) {
        this.getSystemSetting(function (data) {
            let setting = {}

            if (data.logo) {
                setting.logo = data.logo
            }
            if (data.title) {
                setting.title = data.title
            }
            app.$store.commit('setSystemSetting', setting)
        })
    }
}
