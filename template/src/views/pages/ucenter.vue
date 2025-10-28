<template>
    <div>
        <div class="user-container">
            <el-card class="user-profile" shadow="hover" :body-style="{ padding: '0px' }">
                <div class="user-profile-bg" :style="{ backgroundImage: loginConfig('userbg') ? `url(${loginConfig('userbg')})` : '' }"></div>
                <div class="user-avatar-wrap">
                    <el-avatar class="user-avatar" :size="120" :src="currentAvatar" />
                </div>
                <div class="user-info">
                    <div class="info-name">{{ name }}</div>
                </div>
            </el-card>
            <el-card
                class="user-content"
                shadow="hover"
                :body-style="{ padding: '20px 50px', height: '100%', boxSizing: 'border-box' }"
            >
                <el-tabs tab-position="left" v-model="activeName">
                    <el-tab-pane name="label2" label="我的头像" class="user-tabpane">
                        <div class="crop-wrap" v-if="activeName === 'label2'">
                            <vueCropper
                                ref="cropper"
                                :img="imgSrc"
                                :autoCrop="true"
                                :centerBox="true"
                                :full="true"
                                mode="contain"
                                :crossOrigin="false"
                            >
                            </vueCropper>
                        </div>
                        <el-button class="crop-demo-btn" type="primary"
                            >选择图片
                            <input class="crop-input" type="file" name="image" accept="image/*" @change="setImage" />
                        </el-button>
                        <el-button type="success" @click="saveAvatar" :loading="loading" :disabled="!imgSrc || imgSrc === avatarImg">
                            {{ loading ? '上传中...' : '上传并保存' }}
                        </el-button>
                    </el-tab-pane>
                    <el-tab-pane name="label3" label="修改密码" class="user-tabpane">
                        <el-form class="w500" label-position="top">
                            <el-form-item label="旧密码：">
                                <el-input type="password" v-model="form.old"></el-input>
                            </el-form-item>
                            <el-form-item label="新密码：">
                                <el-input type="password" v-model="form.new"></el-input>
                            </el-form-item>
                            <el-form-item label="确认新密码：">
                                <el-input type="password" v-model="form.new1"></el-input>
                            </el-form-item>
                            <el-form-item>
                                <el-button type="primary" @click="onSubmit">保存</el-button>
                            </el-form-item>
                        </el-form>
                    </el-tab-pane>
                </el-tabs>
            </el-card>
        </div>
    </div>
</template>

<script setup lang="ts" name="ucenter">
import service from '@/utils/request';
import { ElLoading, ElMessage } from 'element-plus';
import { computed, onMounted, onUnmounted, reactive, ref } from 'vue';
import { VueCropper } from 'vue-cropper';
import 'vue-cropper/dist/index.css';

const name = localStorage.getItem('username');
const avatarImg = ref("");
avatarImg.value = localStorage.getItem('avatar');

// 当前头像URL（包含时间戳防止缓存）
const currentAvatar = computed(() => {
    if (!avatarImg.value) return '';
    // 添加时间戳参数防止浏览器缓存
    const separator = avatarImg.value.includes('?') ? '&' : '?';
    return `${avatarImg.value}${separator}t=${Date.now()}`;
});

// 配置缓存（响应式）
const configState = ref<Record<string, any>>({});
const getStoredConfig = () => {
    const config = localStorage.getItem('config');
    if (!config) return {} as Record<string, any>;
    try {
        return JSON.parse(config) || {};
    } catch {
        return {};
    }
};
configState.value = getStoredConfig();

const loginConfig = (field: string) => {
    return (configState.value as any)?.[field] ?? '';
};

// 监听配置更新事件
const handleConfigUpdated = (event: CustomEvent) => {
    if (event.detail && event.detail.config) {
        configState.value = event.detail.config;
    }
};

onMounted(() => {
    window.addEventListener('configUpdated', handleConfigUpdated as EventListener);
});

onUnmounted(() => {
    window.removeEventListener('configUpdated', handleConfigUpdated as EventListener);
});

const form = reactive({
    new1: '',
    new: '',
    old: '',
});
const onSubmit = async () => {
    try {
        if (form.new !== form.new1) {
            ElMessage.error('两次输入的密码不一致');
            return;
        }
        var url="/api/api/v1/admin/updatepassword";
        const res = await service.post(url, form);
        if (res.data.code === 200) {
            ElMessage.success(res.data.msg);
        } else {
            ElMessage.error(res.data.msg);
        }
    } catch (e: any) {
        ElMessage.error("未知错误");
    }
};
const activeName = ref('label2');
const imgSrc = ref(localStorage.getItem('avatar'));
const cropper: any = ref();
const loading = ref(false);

const setImage = (e: any) => {
    const file = e.target.files[0];
    if (!file) return;
    
    if (!file.type.includes('image/')) {
        ElMessage.error('请选择图片文件');
        return;
    }
    
    // 检查文件大小（限制为5MB）
    if (file.size > 5 * 1024 * 1024) {
        ElMessage.error('图片大小不能超过5MB');
        return;
    }
    
    const reader = new FileReader();
    reader.onload = (event: any) => {
        imgSrc.value = event.target.result;
    };
    reader.readAsDataURL(file);
};




const saveAvatar = async () => {
    if (!cropper.value) {
        ElMessage.error('请先选择图片');
        return;
    }
    
    if (!imgSrc.value) {
        ElMessage.error('请先选择要裁剪的图片');
        return;
    }
    
    loading.value = true;
    const loadingInstance = ElLoading.service({
        lock: true,
        text: '上传中...',
        background: 'rgba(0, 0, 0, 0.7)'
    });
    
    try {
        // 使用getCropBlob获取更好质量的裁剪图片
        const blob = await new Promise((resolve, reject) => {
            cropper.value.getCropBlob((blob: Blob) => {
                if (blob) {
                    resolve(blob);
                } else {
                    reject(new Error('获取裁剪图片失败'));
                }
            }, 'image/jpeg', 0.9); // 设置为JPEG格式，质量0.9
        }) as Blob;
        
        // 创建FormData
        const formData = new FormData();
        const fileName = `avatar_${Date.now()}.jpg`;
        formData.append('files', blob, fileName);
        
        // 上传头像
        const res = await service.post('/api/api/v1/system/upload', formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        });
        if (res.data.code === 200) {
            var uploadedPath = res.data.data.filePaths[0];
            const res2 = await service.post('/api/api/v1/admin/updateavatar', {avatar: uploadedPath});
            if (res2.data.code === 200) {
                ElMessage.success('头像上传成功');
                // 强制更新响应式数据
                avatarImg.value = uploadedPath;
                imgSrc.value = uploadedPath;
                // 更新localStorage中的头像信息
                localStorage.setItem('avatar', uploadedPath);
                // 刷新页面头像（如果有全局头像显示）
                window.dispatchEvent(new CustomEvent('avatarUpdated', {
                    detail: { avatar: uploadedPath }
                }));
            } else {
                ElMessage.error(res2.data.msg);
            }
        } else {
            ElMessage.error(res.data.msg || '上传失败');
        }
        
    } catch (error: any) {
        if (error.response?.data?.msg) {
            ElMessage.error(error.response.data.msg);
        } else {
            ElMessage.error('上传失败，请重试');
        }
    } finally {
        loading.value = false;
        loadingInstance.close();
    }
};
</script>

<style scoped>
.user-container {
    display: flex;
}

.user-profile {
    position: relative;
}

.user-profile-bg {
    width: 100%;
    height: 200px;
    background-size: cover;
    background-position: center;
    background-repeat: no-repeat;
}

.user-profile {
    width: 500px;
    margin-right: 20px;
    flex: 0 0 auto;
    align-self: flex-start;
}

.user-avatar-wrap {
    position: absolute;
    top: 135px;
    width: 100%;
    text-align: center;
}

.user-avatar {
    border: 5px solid #fff;
    border-radius: 50%;
    overflow: hidden;
    box-shadow: 0 7px 12px 0 rgba(62, 57, 107, 0.16);
}

.user-info {
    text-align: center;
    padding: 80px 0 30px;
}

.info-name {
    margin: 0 0 20px;
    font-size: 22px;
    font-weight: 500;
    color: #373a3c;
}

.info-desc {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 5px;
}

.info-desc,
.info-desc a {
    font-size: 18px;
    color: #55595c;
}

.info-icon {
    margin-top: 10px;
}

.info-icon i {
    font-size: 30px;
    margin: 0 10px;
    color: #343434;
}

.user-content {
    flex: 1;
}

.user-tabpane {
    padding: 10px 20px;
}

.crop-wrap {
    width: 600px;
    height: 350px;
    margin-bottom: 20px;
}

.crop-demo-btn {
    position: relative;
}

.crop-input {
    position: absolute;
    width: 100px;
    height: 40px;
    left: 0;
    top: 0;
    opacity: 0;
    cursor: pointer;
}

.w500 {
    width: 500px;
}

.user-footer {
    display: flex;
    border-top: 1px solid rgba(83, 70, 134, 0.1);
}

.user-footer-item {
    padding: 20px 0;
    width: 33.3333333333%;
    text-align: center;
}

.user-footer > div + div {
    border-left: 1px solid rgba(83, 70, 134, 0.1);
}
</style>

<style>
.el-tabs.el-tabs--left {
    height: 100%;
}
</style>
