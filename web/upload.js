document.getElementById('uploadButton').addEventListener('click', async () => {
    const fileInput = document.getElementById('fileInput');
    const statusText = document.getElementById('statusText');

    if (fileInput.files.length === 0) {
        statusText.innerText = '请选择一个文件上传';
        return;
    }

    const file = fileInput.files[0];
    statusText.innerText = '获取签名URL...';




    try {
        // 请求服务端获取预签名的URL
        const response = await fetch('http://localhost:8080/api/v1/assets/aliyun/oss/policy', {
            method: "POST",
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDU5ODA4OTQsImlzcyI6Im1lbG9uIiwiVXNlcklEIjoiZTBmMWNiODZjN2U2NDQ2Y2E5MjFhMGUwYTA4YzcxOGIifQ.rbTkxYvcEI3Mgx2_1I4-zGs4GmLoYKl5z420zT1IZ1M`
            },
            body: JSON.stringify({
              "storage": "users",
              "ext": "png"
            })
        });
        if (!response.ok) {
            throw new Error('无法获取签名URL');
        }

        const resp = await response.json();
        const data = resp.data;
        const formData = new FormData();
        formData.append('key', data.cos_path);
        formData.append('policy', data.policy);
        formData.append('OSSAccessKeyId', data.accessKey_id);
        formData.append('signature', data.signature);
        formData.append('file', file);

        // 上传到OSS
        fetch(`${data.upload_url}`, {
            method: 'POST',
            body: formData
        })
        .then(response => {
            if (response.ok) {
                alert('文件上传成功, ' + data.asset_url);
            } else {
                alert('文件上传失败');
            }
        })
        .catch(error => {
            console.error('上传错误:', error);
        });
    } catch (error) {
        statusText.innerText = `错误: ${error.message}`;
    }
});