### gokeeper

------------
gokeeper 是个尚处于密集开发阶段的微服务框架(或者称为微服务管理框架也行)。它目前是在复刻[dmicro](https://dmicro.vprix.com/#/dserver/quickstart)项目，并做一些积极的优化。

### gokeeper设计初衷
dmicro是[刘志铭](https://github.com/osgochina)大佬写的一款微服务框架，代码结构清晰，文档和注释都比较完善。个人觉得其中比较有趣的组件有进程管理、服务平滑重启、交互式shell、rpc组件等。目前dmicro还在早期开发阶段，也在持续更新。

创建gokeeper的初衷主要是针对dmicro存在问题做一系列优化：
- 设计更清晰；gokeeper将致力于让代码更符合go语言习惯，尽量减少冗余，让概念不易混淆；
- 模块尽量解耦，减少相互调用和反复横跳的情况；例如交互式shell的服务端采用gin + unixDomainSocket 的方式实现，不依赖内部的drpc；
- 更好的交互式shell体验，提供更精细的控制；
- 尽量避免过度设计，例如dmicro中的drpc模块；
- 提高框架兼容性，例如让rpcx、erpc等rpc框架能够兼容平滑重启，这样就没必要重复造轮子，同时也能让已经使用了这些rpc框架的小伙伴能快速迁移；

总之，就是让项目逻辑更清晰，更便于修改；同时增加兼容性；并在此基础上增加更多好用的功能。

### gokeeper todo
  - [ ] 重新设计dmicro的除drpc之外的所有组件
  - [ ] 兼容一些常见的rpc框架
  - [ ] 其他dmicro尚在计划中的功能
  - [ ] more
