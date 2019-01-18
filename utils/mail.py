# coding: utf-8

import smtplib
from email.mime.multipart import MIMEMultipart
from email.mime.text import MIMEText
from email.header import Header

class SendMail:

    smtpserver = 'smtp.exmail.qq.com'
    username = 'ai-ops@jollycorp.com'
    password = 'Www263com'
    sender = 'ai-ops@jollycorp.com'
    receiver = 'op@jollycorp.com'
    subject = '服务器到期提醒'
    subject=Header(subject, 'utf-8').encode()

    @classmethod
    def send(cls,content):
        msg = MIMEMultipart('mixed')
        msg['Subject'] = cls.subject
        msg['From'] = cls.sender
        #msg['To'] = ";".join(receiver)
        msg['To'] = cls.receiver

        # 邮件内容
        text = content
        text_plain = MIMEText(text, 'plain', 'utf-8')
        msg.attach(text_plain)

        # 发送邮件
        smtp = smtplib.SMTP()
        smtp.connect(cls.smtpserver)
        smtp.login(cls.username, cls.password)
        smtp.sendmail(cls.sender, cls.receiver, msg.as_string())
        smtp.quit()
