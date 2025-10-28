<template>
    <div class="login-bg" :style="{ 
            backgroundImage:  getConfig('loginbg') ? `url(${ getConfig('loginbg')})` : '',
            backgroundPosition: 'center',
            backgroundSize: 'cover',
            backgroundRepeat: 'no-repeat'
        }">
        <div class="login-container">
            <div class="login-header">
                <img class="logo mr10" :src=" getConfig('logo')" alt="" />
                <div class="login-title">{{  getConfig('title') }}</div>
            </div>
            <el-form :model="param" :rules="rules" ref="login" size="large">
                <el-form-item prop="username">  
                    <el-input v-model="param.username" placeholder="用户名">
                        <template #prepend>
                            <el-icon>
                                <User />
                            </el-icon>
                        </template>
                    </el-input>
                </el-form-item>
                <el-form-item prop="password">
                    <el-input
                        type="password"
                        placeholder="密码"
                        v-model="param.password"
                        @keyup.enter="submitForm(login)"
                    >
                        <template #prepend>
                            <el-icon>
                                <Lock />
                            </el-icon>
                        </template>
                    </el-input>
                </el-form-item>
                <!-- <div class="pwd-tips">
                    <el-checkbox class="pwd-checkbox" v-model="checked" label="记住密码" />
                    <el-link type="primary" @click="$router.push('/reset-pwd')">忘记密码</el-link>
                </div> -->
                <el-button class="login-btn" type="primary" size="large" @click="submitForm(login)">登录</el-button>
                <!-- <p class="login-tips">Tips : 用户名和密码随便填。</p>
                <p class="login-text">
                    没有账号？<el-link type="primary" @click="$router.push('/register')">立即注册</el-link>
                </p> -->
            </el-form>
        </div>
    </div>
</template>

<script setup lang="ts">
import { usePermissStore } from '@/store/permiss';
import { useTabsStore } from '@/store/tabs';
import axios from 'axios';
import type { FormInstance, FormRules } from 'element-plus';
import { ElMessage } from 'element-plus';
import { onMounted, reactive, ref } from 'vue';
import { useRouter } from 'vue-router';

interface LoginInfo {
    username: string;
    password: string;
}

const lgStr = localStorage.getItem('login-param');
const defParam = lgStr ? JSON.parse(lgStr) : null;
const checked = ref(lgStr ? true : false);
// 配置缓存（同步给模板使用）
const configMap = ref<Record<string, any>>({});

// 同步获取字段值，供模板直接使用
const getConfig = (field: string): string => {
    const val = configMap.value?.[field];
    return val == null ? '' : String(val);
};

// 加载配置：优先本地缓存，没有则请求接口后写入本地
const loadConfig = async () => {
    const cached = localStorage.getItem('config');
    let data: any | null = null;

    if (cached) {
        try {
            data = JSON.parse(cached);
        } catch (e) {
            localStorage.removeItem('config');
        }
    }
    if (!data) {
        const res = await axios.post('/api/api/v1/config/getconfig');
        if (res?.data?.code === 200) {
            data = res.data.data || {};
            localStorage.setItem('config', JSON.stringify(data));
        } else {
            data = {};
        }
    }

    configMap.value = data || {};
};

onMounted(() => {
    loadConfig();
});



const router = useRouter();
const param = reactive<LoginInfo>({
    username: defParam ? defParam.username : '',
    password: defParam ? defParam.password : '',
});

const rules: FormRules = {
    username: [
        {
            required: true,
            message: '请输入用户名',
            trigger: 'blur',
        },
    ],
    password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
};
const permiss = usePermissStore();
const login = ref<FormInstance>();
const submitForm = async (formEl: FormInstance | undefined) => {
    if (!formEl) return;
    formEl.validate(async (valid: boolean) => {
        if (!valid) {
            ElMessage.error('登录失败');
            return false;
        }
        try {
            const res = await axios.post('/api/api/v1/login', {
                username: param.username,
                password: param.password
            });
            // 登录成功
            if (res.data.code == 200) {
                ElMessage.success('登录成功');
                localStorage.setItem('username', param.username);
                localStorage.setItem('avatar', res.data.data.admin.avatar);
                localStorage.setItem('token', res.data.data.token);
                localStorage.setItem('rulearr', res.data.data.rulearr);
                permiss.handleSet(res.data.data.rulearr);
                router.push('/');
                if (checked.value) {
                    localStorage.setItem('login-param', JSON.stringify(param));
                } else {
                    localStorage.removeItem('login-param');
                }
            }else{
                ElMessage.error(res.data.msg);
                return false;
            }
        } catch (error: any) {
            // 这里打印后端返回的详细内容
            if (error.response) {
                ElMessage.error(error.response.data?.msg || '登录失败');
            } else {
                ElMessage.error('请求异常');
            }
        }
        return false;
    });
};

const tabs = useTabsStore();
tabs.clearTabs();
</script>

<style scoped>
.login-bg {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100vh;

}

.login-header {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 40px;
}

.logo {
    width: 35px;
}

.login-title {
    font-size: 22px;
    color: #333;
    font-weight: bold;
}

.login-container {
    width: 450px;
    border-radius: 5px;
    background: #fff;
    padding: 40px 50px 50px;
    box-sizing: border-box;
}

.pwd-tips {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 14px;
    margin: -10px 0 10px;
    color: #787878;
}

.pwd-checkbox {
    height: auto;
}

.login-btn {
    display: block;
    width: 100%;
}

.login-tips {
    font-size: 12px;
    color: #999;
}

.login-text {
    display: flex;
    align-items: center;
    margin-top: 20px;
    font-size: 14px;
    color: #787878;
}
</style>
