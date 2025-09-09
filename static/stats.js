/**
 * Webbleen 博客统计脚本
 * 用于在 Hugo 博客中集成访问统计功能
 */

(function() {
    'use strict';
    
    // 配置
    const CONFIG = {
        apiUrl: 'http://localhost:8000/api', // 修改为你的API地址
        sessionKey: 'webbleen_session_id',
        visitKey: 'webbleen_visit_recorded'
    };
    
    // 获取或生成会话ID
    function getSessionId() {
        let sessionId = localStorage.getItem(CONFIG.sessionKey);
        if (!sessionId) {
            sessionId = 'session_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
            localStorage.setItem(CONFIG.sessionKey, sessionId);
        }
        return sessionId;
    }
    
    // 检测设备类型
    function getDeviceType() {
        const userAgent = navigator.userAgent;
        if (/Mobile|Android|iPhone|iPad/.test(userAgent)) {
            return 'mobile';
        } else if (/Tablet|iPad/.test(userAgent)) {
            return 'tablet';
        } else {
            return 'desktop';
        }
    }
    
    // 检测浏览器
    function getBrowser() {
        const userAgent = navigator.userAgent;
        if (userAgent.indexOf('Chrome') > -1) return 'Chrome';
        if (userAgent.indexOf('Firefox') > -1) return 'Firefox';
        if (userAgent.indexOf('Safari') > -1) return 'Safari';
        if (userAgent.indexOf('Edge') > -1) return 'Edge';
        if (userAgent.indexOf('Opera') > -1) return 'Opera';
        return 'Unknown';
    }
    
    // 检测操作系统
    function getOS() {
        const userAgent = navigator.userAgent;
        if (userAgent.indexOf('Windows') > -1) return 'Windows';
        if (userAgent.indexOf('Mac') > -1) return 'macOS';
        if (userAgent.indexOf('Linux') > -1) return 'Linux';
        if (userAgent.indexOf('Android') > -1) return 'Android';
        if (userAgent.indexOf('iOS') > -1) return 'iOS';
        return 'Unknown';
    }
    
    // 获取地理位置信息（简化版）
    function getLocation() {
        // 这里可以集成第三方地理位置API
        return {
            country: 'Unknown',
            city: 'Unknown'
        };
    }
    
    // 记录访问
    function recordVisit() {
        // 检查是否已经记录过本次访问
        const visitKey = CONFIG.visitKey + '_' + window.location.pathname;
        if (sessionStorage.getItem(visitKey)) {
            return;
        }
        
        const visitData = {
            page: window.location.pathname,
            session_id: getSessionId(),
            device: getDeviceType(),
            browser: getBrowser(),
            os: getOS(),
            country: getLocation().country,
            city: getLocation().city
        };
        
        // 发送统计数据
        fetch(CONFIG.apiUrl + '/visit', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(visitData)
        }).then(response => {
            if (response.ok) {
                // 标记已记录
                sessionStorage.setItem(visitKey, 'true');
                console.log('Visit recorded successfully');
            }
        }).catch(error => {
            console.error('Error recording visit:', error);
        });
    }
    
    // 获取统计数据
    function getStats() {
        return fetch(CONFIG.apiUrl + '/stats/visits')
            .then(response => response.json())
            .then(data => {
                if (data.code === 200) {
                    return data.data;
                }
                throw new Error(data.msg);
            });
    }
    
    // 获取内容统计
    function getContentStats() {
        return fetch(CONFIG.apiUrl + '/stats/content')
            .then(response => response.json())
            .then(data => {
                if (data.code === 200) {
                    return data.data;
                }
                throw new Error(data.msg);
            });
    }
    
    // 获取热门页面
    function getTopPages(limit = 10) {
        return fetch(CONFIG.apiUrl + '/stats/pages?limit=' + limit)
            .then(response => response.json())
            .then(data => {
                if (data.code === 200) {
                    return data.data;
                }
                throw new Error(data.msg);
            });
    }
    
    // 获取访问趋势
    function getVisitTrend(days = 30) {
        return fetch(CONFIG.apiUrl + '/stats/trend?days=' + days)
            .then(response => response.json())
            .then(data => {
                if (data.code === 200) {
                    return data.data;
                }
                throw new Error(data.msg);
            });
    }
    
    // 获取用户行为分析
    function getUserBehavior() {
        return fetch(CONFIG.apiUrl + '/stats/behavior')
            .then(response => response.json())
            .then(data => {
                if (data.code === 200) {
                    return data.data;
                }
                throw new Error(data.msg);
            });
    }
    
    // 页面加载完成后记录访问
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', recordVisit);
    } else {
        recordVisit();
    }
    
    // 页面卸载时记录停留时间
    window.addEventListener('beforeunload', function() {
        const startTime = performance.timing.navigationStart;
        const endTime = Date.now();
        const stayTime = Math.round((endTime - startTime) / 1000);
        
        // 可以发送停留时间数据
        fetch(CONFIG.apiUrl + '/visit', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                page: window.location.pathname,
                session_id: getSessionId(),
                stay_time: stayTime
            })
        }).catch(error => {
            console.error('Error recording stay time:', error);
        });
    });
    
    // 暴露全局API
    window.WebbleenStats = {
        recordVisit: recordVisit,
        getStats: getStats,
        getContentStats: getContentStats,
        getTopPages: getTopPages,
        getVisitTrend: getVisitTrend,
        getUserBehavior: getUserBehavior,
        getSessionId: getSessionId
    };
    
})();
