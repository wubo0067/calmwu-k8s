/*
 * @Author: calm.wu
 * @Date: 2020-09-02 17:34:43
 * @Last Modified by: calm.wu
 * @Last Modified time: 2020-09-02 18:00:24
 */

package tinycontroller

import "errors"

var ErrResourceNotSupport = errors.New("resource not support")

var ErrCacheSyncTimeout = errors.New("cache resource sync timeout")
