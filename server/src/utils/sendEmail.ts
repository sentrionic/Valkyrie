import * as nodemailer from 'nodemailer';

// async..await is not allowed in global scope, must use a wrapper
export async function sendEmail(to: string, html: string): Promise<void> {
  // create Nodemailer SES transporter
  const transporter = nodemailer.createTransport({
    port: 587,
    service: 'gmail',
    secure: true,
    auth: {
      user: process.env.GMAIL_USER,
      pass: process.env.GMAIL_PASSWORD,
    },
    debug: true,
  });

  // send mail with defined transport object
  await transporter.sendMail({
    from: '"Valkyrie Team" <support@vakyrie.com>', // sender address
    to: to, // list of receivers
    subject: 'Change Password', // Subject line
    html,
  });
}
