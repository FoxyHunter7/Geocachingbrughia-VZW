<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class MessageLanguage extends Model
{
    use HasFactory;

    protected $fillable = [
        'message_id',
        'lang_code',
        'title',
        'body'
    ];

    protected $hidden = [
        'id',
        'message_id',
        'created_at',
        'updated_at'
    ];
}
