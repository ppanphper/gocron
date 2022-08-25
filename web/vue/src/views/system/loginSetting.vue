<template>
  <el-card>
    <el-tabs v-model="activeName">
      <el-tab-pane label="Ldap登录认证" name="ldap">
        <el-form ref="ldapForm" :model="ldapSetting" :rules="ldapSettingRules" label-width="120px"
                 style="width: 700px;padding-top: 15px">
          <el-row>
            <el-col :span="6">
              <el-form-item label="Enable" prop="enable">
                <el-switch active-value="1" inactive-value="0" v-model="ldapSetting.enable"></el-switch>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="Url" prop="url" required>
                <el-input v-model="ldapSetting.url"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="绑定DN" prop="bind_dn" required>
                <el-input v-model="ldapSetting.bind_dn"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="绑定密码" prop="bind_password" required>
                <el-input type="password" v-model="ldapSetting.bind_password"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="筛选范围" prop="base_dn" required>
                <el-input v-model.number="ldapSetting.base_dn"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="筛选规则" prop="filter_rule" required>
                <el-input v-model.number="ldapSetting.filter_rule" placeholder="筛选规则,如：(&(cn=%s))"></el-input>
              </el-form-item>
            </el-col>
            <el-col :span="24">
              <el-form-item label="LDAP邮箱属性" prop="ldap_email_attribute" required>
                <el-input v-model.number="ldapSetting.ldap_email_attribute"></el-input>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item>
            <el-button @click="ldapConnectionTest">LDAP连接测试</el-button>
            <el-button type="primary" @click="ldapSettingSubmit()">提交</el-button>
          </el-form-item>
        </el-form>

        <el-dialog v-model="ldapTestVisible" title="LDAP Test">
          <el-form ref="ldapTestForm" :model="testForm" :rules="testFormRules" label-width="120px">
            <el-form-item label="用户名" prop="username" required>
              <el-input v-model="testForm.username" autocomplete="off"/>
            </el-form-item>
            <el-form-item label="密码" prop="password" required>
              <el-input type="password" v-model="testForm.password" autocomplete="off"/>
            </el-form-item>
          </el-form>
          <template #footer>
            <span class="dialog-footer">
              <el-button @click="ldapTestVisible = false">取消</el-button>
              <el-button type="primary" @click="ldapConnectionTestSubmit">测试连接</el-button>
            </span>
          </template>
        </el-dialog>
      </el-tab-pane>
      <!--        <el-tab-pane label="oauth2.0登录认证" name="sso">
                sso
              </el-tab-pane>-->
    </el-tabs>
  </el-card>
</template>

<script>
import httpClient from '../../utils/httpClient'

export default {
  name: 'login-setting',
  data() {
    return {
      activeName: 'ldap',
      ldapTestVisible: false,
      ldapSetting: {
        enable: '0',
        url: '',
        bind_dn: '',
        bind_password: '',
        ldap_email_attribute: 'mail',
        base_dn: 'ou=users,dc=example,dc=com',
        filter_rule: ''
      },
      ldapSettingRules: {
        url: [
          {required: true, message: 'Url不能为空', trigger: 'blur'}
        ],
        bind_dn: [
          {required: true, message: '绑定DN不能为空', trigger: 'blur'}
        ],
        bind_password: [
          {required: true, message: '绑定密码不能为空', trigger: 'blur'}
        ],
        base_dn: [
          {required: true, message: '筛选范围不能为空', trigger: 'blur'}
        ],
        filter_rule: [
          {required: true, message: '筛选规则不能为空', trigger: 'blur'}
        ],
        ldap_email_attribute: [
          {required: true, message: 'LDAP邮箱属性不能为空', trigger: 'blur'}
        ],
      },
      testForm: {},
      testFormRules: {
        username: [
          {required: true, message: '用户名不能为空', trigger: 'blur'}
        ],
        password: [
          {required: true, message: '密码不能为空', trigger: 'blur'}
        ]
      },
    }
  },
  mounted() {
    this.renderSetting()
  },
  methods: {
    ldapSettingSubmit() {
      let _this = this
      this.$refs['ldapForm'].validate(valid => {
        if (!valid) {
          return false
        }
        httpClient.post('/system/ldap/update', this.ldapSetting, function () {
          _this.$message.success('设置成功')
        })
        console.log(this.ldapSetting)
      })
    },
    ldapConnectionTest() {
      let _this = this
      this.$refs['ldapForm'].validate(valid => {
        if (!valid) {
          return false
        }
        _this.ldapTestVisible = true
      })
    },
    ldapConnectionTestSubmit() {
      let _this = this
      _this.$refs['ldapTestForm'].validate(valid => {
        if (!valid) {
          return false
        }
        let data = Object.assign({}, this.ldapSetting, this.testForm)
        httpClient.post("/system/ldap/test", data, function () {
          _this.$message.success("登录验证成功,请保存配置")
        })
      })
    },
    renderSetting() {
      let _this = this
      httpClient.get('/system/ldap', {}, function (data) {
        Object.assign(_this.ldapSetting, data)
      })
    }
  }
}
</script>
